package graph

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/warerastats/api/graph/model"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ---------------- DamageRanking / WageRanking / OrderBook ----------------

func (r *damageRankingResolver) User(ctx context.Context, obj *model.DamageRanking) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *wageRankingResolver) User(ctx context.Context, obj *model.WageRanking) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *orderBookResolver) Bids(ctx context.Context, obj *model.OrderBook) ([]*model.OrderBookLevel, error) {
	rows, err := r.Colls.Trackers.TradeOffer.GetOpenOffers(ctx, obj.ItemCode, "BUY", true)
	if err != nil {
		return nil, err
	}
	out := make([]*model.OrderBookLevel, len(rows))
	for i, l := range rows {
		out[i] = &model.OrderBookLevel{Price: l.Price, Remaining: int32(l.Remaining)}
	}
	return out, nil
}

func (r *orderBookResolver) Asks(ctx context.Context, obj *model.OrderBook) ([]*model.OrderBookLevel, error) {
	rows, err := r.Colls.Trackers.TradeOffer.GetOpenOffers(ctx, obj.ItemCode, "SELL", false)
	if err != nil {
		return nil, err
	}
	out := make([]*model.OrderBookLevel, len(rows))
	for i, l := range rows {
		out[i] = &model.OrderBookLevel{Price: l.Price, Remaining: int32(l.Remaining)}
	}
	return out, nil
}

// ---------------- Query: get by id ----------------

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return loadUser(ctx, id)
}

func (r *queryResolver) Country(ctx context.Context, id string) (*model.Country, error) {
	return loadCountry(ctx, id)
}

func (r *queryResolver) Party(ctx context.Context, id string) (*model.Party, error) {
	return loadParty(ctx, id)
}

func (r *queryResolver) Mu(ctx context.Context, id string) (*model.Mu, error) {
	return loadMu(ctx, id)
}

func (r *queryResolver) Region(ctx context.Context, id string) (*model.Region, error) {
	return loadRegion(ctx, id)
}

func (r *queryResolver) Battle(ctx context.Context, id string) (*model.Battle, error) {
	return loadBattle(ctx, id)
}

func (r *queryResolver) Item(ctx context.Context, id string) (*model.Item, error) {
	return loadItem(ctx, id)
}

func (r *queryResolver) Company(ctx context.Context, id string) (*model.Company, error) {
	return loadCompany(ctx, id)
}

func (r *queryResolver) TradeOffer(ctx context.Context, id string) (*model.TradeOffer, error) {
	oid, err := oidOf(id)
	if err != nil {
		return nil, err
	}
	o, err := r.Colls.Trackers.TradeOffer.Get(ctx, oid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return toTradeOffer(o), nil
}

// ---------------- Query: search ----------------

func (r *queryResolver) Search(ctx context.Context, term string, limit *int32) ([]model.SearchResult, error) {
	n := limitOf(limit, 10)
	var out []model.SearchResult

	users, err := r.Colls.Trackers.User.Search(ctx, term, n)
	if err != nil {
		return nil, err
	}
	for i := range users {
		out = append(out, toUser(&users[i]))
	}
	parties, err := r.Colls.Trackers.Party.Search(ctx, term, n)
	if err != nil {
		return nil, err
	}
	for i := range parties {
		out = append(out, toParty(&parties[i]))
	}
	countries, err := r.Colls.Trackers.Country.Search(ctx, term, n)
	if err != nil {
		return nil, err
	}
	for i := range countries {
		out = append(out, toCountry(&countries[i]))
	}
	mus, err := r.Colls.Trackers.Mu.Search(ctx, term, n)
	if err != nil {
		return nil, err
	}
	for i := range mus {
		out = append(out, toMu(&mus[i]))
	}
	return out, nil
}

// ---------------- Query: root lists ----------------

func (r *queryResolver) Battles(ctx context.Context, first *int32, after *string, filter *model.BattleFilter) ([]*model.Battle, error) {
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	rows, err := r.Colls.Trackers.Battle.List(ctx, toBattleFilter(filter), before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toBattle), nil
}

func (r *queryResolver) Countries(ctx context.Context) ([]*model.Country, error) {
	rows, err := r.Colls.Trackers.Country.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toCountry), nil
}

func (r *queryResolver) Regions(ctx context.Context) ([]*model.Region, error) {
	rows, err := r.Colls.Trackers.Region.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toRegion), nil
}

func (r *queryResolver) Parties(ctx context.Context, first *int32, after *string, activeOnly *bool) ([]*model.Party, error) {
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	active := activeOnly == nil || *activeOnly
	rows, err := r.Colls.Trackers.Party.ListActive(ctx, active, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toParty), nil
}

func (r *queryResolver) Mus(ctx context.Context, first *int32, after *string, activeOnly *bool) ([]*model.Mu, error) {
	before, err := cursorPtr(after)
	if err != nil {
		return nil, err
	}
	active := activeOnly == nil || *activeOnly
	rows, err := r.Colls.Trackers.Mu.ListActive(ctx, active, before, limitOf(first, 20))
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toMu), nil
}

func (r *queryResolver) OrderBook(ctx context.Context, itemCode string) (*model.OrderBook, error) {
	return &model.OrderBook{ItemCode: itemCode}, nil
}

// ---------------- Query: time-series / derived ----------------

func (r *queryResolver) ItemCandles(ctx context.Context, itemCode string, from time.Time, to time.Time) ([]*model.ItemCandle, error) {
	if err := enforceTimeWindow(ctx, from, to, reportTimeWindowDays); err != nil {
		return nil, err
	}
	rows, err := r.Colls.Processed.Candles.ItemCandle.GetRange(ctx, itemCode, from, to)
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toItemCandle), nil
}

func (r *queryResolver) WageCandles(ctx context.Context, from time.Time, to time.Time) ([]*model.WageCandle, error) {
	if err := enforceTimeWindow(ctx, from, to, reportTimeWindowDays); err != nil {
		return nil, err
	}
	rows, err := r.Colls.Processed.Candles.WageCandle.GetRange(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toWageCandle), nil
}

func (r *queryResolver) MarketStates(ctx context.Context, from time.Time, to time.Time) ([]*model.MarketState, error) {
	if err := enforceTimeWindow(ctx, from, to, reportTimeWindowDays); err != nil {
		return nil, err
	}
	rows, err := r.Colls.Processed.Reports.MarketState.GetRange(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toMarketState), nil
}

func (r *queryResolver) LatestMarketState(ctx context.Context) (*model.MarketState, error) {
	st, ok, err := r.Colls.Processed.Reports.MarketState.GetLatest(ctx)
	if err != nil || !ok {
		return nil, err
	}
	return toMarketState(*st), nil
}

func (r *queryResolver) WageMarketStates(ctx context.Context, from time.Time, to time.Time) ([]*model.WageMarketState, error) {
	if err := enforceTimeWindow(ctx, from, to, reportTimeWindowDays); err != nil {
		return nil, err
	}
	rows, err := r.Colls.Processed.Reports.WageMarketState.GetRange(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toWageMarketState), nil
}

func (r *queryResolver) LatestWageMarketState(ctx context.Context) (*model.WageMarketState, error) {
	st, ok, err := r.Colls.Processed.Reports.WageMarketState.GetLatest(ctx)
	if err != nil || !ok {
		return nil, err
	}
	return toWageMarketState(*st), nil
}

func (r *queryResolver) Inflation(ctx context.Context, from time.Time, to time.Time) ([]*model.InflationPoint, error) {
	if err := enforceTimeWindow(ctx, from, to, reportTimeWindowDays); err != nil {
		return nil, err
	}
	rows, err := r.Colls.Processed.Estimators.Inflation.GetRange(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toInflationPoint), nil
}

func (r *queryResolver) DismantleReports(ctx context.Context, from time.Time, to time.Time) ([]*model.DismantleReport, error) {
	if err := enforceTimeWindow(ctx, from, to, reportTimeWindowDays); err != nil {
		return nil, err
	}
	rows, err := r.Colls.Processed.Reports.DismantleReport.GetRange(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toDismantleReport), nil
}

func (r *queryResolver) ItemMarketReport(ctx context.Context, itemCode string) (*model.ItemMarketReport, error) {
	rep, ok, err := r.Colls.Processed.Reports.ItemMarketReport.Get(ctx, itemCode)
	if err != nil || !ok {
		return nil, err
	}
	return toItemMarketReport(rep), nil
}

func (r *queryResolver) CasesReports(ctx context.Context) ([]*model.CasesReport, error) {
	rows, err := r.Colls.Processed.Reports.CasesReport.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapVal(rows, toCasesReport), nil
}

func (r *queryResolver) EquipmentPricing(ctx context.Context, itemCode string, windowDays int32) (*model.EquipmentPricing, error) {
	win, ok, err := r.Colls.Processed.Reports.EquipmentPricing.GetWindow(ctx, itemCode, int(windowDays))
	if err != nil {
		return nil, err
	}
	skills, err := r.Colls.Processed.Reports.EquipmentPricing.GetSkills(ctx, itemCode, int(windowDays))
	if err != nil {
		return nil, err
	}

	out := &model.EquipmentPricing{ItemCode: itemCode, WindowDays: windowDays}
	if ok {
		out.Window = &model.EquipmentWindowPrice{
			WeightedAvg: win.WeightedAvg, Volume: int32(win.Volume), Count: int32(win.Count), UpdatedAt: win.UpdatedAt,
		}
	}
	out.Skills = make([]*model.EquipmentSkillPrice, len(skills))
	for i, s := range skills {
		out.Skills[i] = &model.EquipmentSkillPrice{
			SkillKey: s.SkillKey, Skills: floatEntries(s.Skills), Min: s.Min, Max: s.Max,
			Avg: s.Avg, Volume: int32(s.Volume), UpdatedAt: s.UpdatedAt,
		}
	}
	return out, nil
}

// ---------------- Query: leaderboards ----------------

func (r *queryResolver) TopDamage(ctx context.Context, from time.Time, to time.Time, limit *int32) ([]*model.DamageRanking, error) {
	if err := enforceTimeWindow(ctx, from, to, reportTimeWindowDays); err != nil {
		return nil, err
	}
	n := limitOf(limit, 10)
	rows, err := r.Colls.Trackers.Damage.AggregateUserDamage(ctx, from, to)
	if err != nil {
		return nil, err
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].TotalDamage > rows[j].TotalDamage })
	if len(rows) > n {
		rows = rows[:n]
	}
	out := make([]*model.DamageRanking, len(rows))
	for i, a := range rows {
		out[i] = &model.DamageRanking{UserID: a.UserID.Hex(), TotalDamage: a.TotalDamage, BattleCount: int32(a.BattleCount)}
	}
	return out, nil
}

func (r *queryResolver) TopWageEarners(ctx context.Context, from time.Time, to time.Time, limit *int32) ([]*model.WageRanking, error) {
	if err := enforceTimeWindow(ctx, from, to, reportTimeWindowDays); err != nil {
		return nil, err
	}
	rows, err := r.Colls.Transactions.WageTransaction.TopPaidEmployees(ctx, from, to, limitOf(limit, 10), false)
	if err != nil {
		return nil, err
	}
	out := make([]*model.WageRanking, len(rows))
	for i, t := range rows {
		out[i] = &model.WageRanking{UserID: t.ID.Hex(), Total: t.Total, Count: int32(t.Count)}
	}
	return out, nil
}

func (r *queryResolver) TopWagePayers(ctx context.Context, from time.Time, to time.Time, limit *int32) ([]*model.WageRanking, error) {
	if err := enforceTimeWindow(ctx, from, to, reportTimeWindowDays); err != nil {
		return nil, err
	}
	n := limitOf(limit, 10)
	rows, err := r.Colls.Transactions.WageTransaction.PaidByEmployer(ctx, from, to)
	if err != nil {
		return nil, err
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Total > rows[j].Total })
	if len(rows) > n {
		rows = rows[:n]
	}
	out := make([]*model.WageRanking, len(rows))
	for i, t := range rows {
		out[i] = &model.WageRanking{UserID: t.ID.Hex(), Total: t.Total, Count: int32(t.Count)}
	}
	return out, nil
}

func (r *queryResolver) TopFlippers(ctx context.Context, limit *int32) ([]*model.UserFlipState, error) {
	rows, err := r.Colls.Processed.Estimators.UserFlipState.TopByProfit(ctx, limitOf(limit, 10))
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toUserFlipState), nil
}

// DamageRanking returns DamageRankingResolver implementation.
func (r *Resolver) DamageRanking() DamageRankingResolver { return &damageRankingResolver{r} }

// OrderBook returns OrderBookResolver implementation.
func (r *Resolver) OrderBook() OrderBookResolver { return &orderBookResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// WageRanking returns WageRankingResolver implementation.
func (r *Resolver) WageRanking() WageRankingResolver { return &wageRankingResolver{r} }

type damageRankingResolver struct{ *Resolver }
type orderBookResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type wageRankingResolver struct{ *Resolver }

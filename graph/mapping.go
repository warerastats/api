package graph

import (
	"time"

	"github.com/warerastats/api/graph/model"
	"github.com/warerastats/models/models/enums"
	"github.com/warerastats/models/models/stores/events"
	processedcandles "github.com/warerastats/models/models/stores/processed/candles"
	processedestimators "github.com/warerastats/models/models/stores/processed/estimators"
	processedreports "github.com/warerastats/models/models/stores/processed/reports"
	"github.com/warerastats/models/models/stores/trackers"
	"github.com/warerastats/models/models/stores/transactions"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ---------- id / scalar helpers ----------

func hexPtr(id *bson.ObjectID) *string {
	if id == nil {
		return nil
	}
	s := id.Hex()
	return &s
}

// atOf derives an event/document timestamp from its ObjectID.
func atOf(id bson.ObjectID) time.Time { return id.Timestamp() }

// cursorPtr parses an optional hex `after` argument into an *ObjectID cursor.
func cursorPtr(after *string) (*bson.ObjectID, error) {
	if after == nil || *after == "" {
		return nil, nil
	}
	id, err := bson.ObjectIDFromHex(*after)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// limitOf clamps an optional `first` argument to a sane page size.
func limitOf(first *int32, def int) int {
	if first == nil || *first <= 0 {
		return def
	}
	if *first > 200 {
		return 200
	}
	return int(*first)
}

func floatEntries(m map[string]float64) []*model.FloatEntry {
	out := make([]*model.FloatEntry, 0, len(m))
	for k, v := range m {
		out = append(out, &model.FloatEntry{Key: k, Value: v})
	}
	return out
}

func intEntries(m map[string]int) []*model.IntEntry {
	out := make([]*model.IntEntry, 0, len(m))
	for k, v := range m {
		out = append(out, &model.IntEntry{Key: k, Value: int32(v)})
	}
	return out
}

// ---------- tracker mappers ----------

func toUser(u *trackers.User) *model.User {
	if u == nil {
		return nil
	}
	return &model.User{
		ID:           u.ID.Hex(),
		Username:     u.Username,
		Level:        int32(u.Level),
		AvatarURL:    u.AvatarUrl,
		MilitaryRank: int32(u.MilitaryRank),
		LastDate:     u.LastDate,
		LastSeen:     u.LastSeen,
		CountryID:    u.CountryID.Hex(),
		CompanyID:    hexPtr(u.CompanyID),
		PartyID:      hexPtr(u.PartyID),
		MuID:         hexPtr(u.MuID),
		WealthMap:    u.Wealth,
		SkillsMap:    u.Skills,
	}
}

func toCountry(c *trackers.Country) *model.Country {
	if c == nil {
		return nil
	}
	return &model.Country{
		ID:    c.ID.Hex(),
		Name:  c.Name,
		Code:  c.Code,
		Money: c.Money,
		Taxes: &model.Taxes{
			Income:   c.Taxes.Income,
			Market:   c.Taxes.Market,
			SelfWork: c.Taxes.SelfWork,
		},
		Specialisation: c.SpecialisationItemCode,
		RulingPartyID:  hexPtr(c.RulingPartyID),
		AllianceID:     hexPtr(c.AllianceID),
	}
}

func toRegion(r *trackers.Region) *model.Region {
	if r == nil {
		return nil
	}
	neighbors := make([]string, len(r.NeighborRegionIDs))
	for i, id := range r.NeighborRegionIDs {
		neighbors[i] = id.Hex()
	}
	return &model.Region{
		ID:                r.ID.Hex(),
		Name:              r.Name,
		IsCapital:         r.IsCapital,
		IsLinkedToCapital: r.IsLinkedToCapital,
		Resistance:        r.Resistance,
		MaxResistance:     r.MaxResistance,
		CountryID:         r.CountryID.Hex(),
		InitialCountryID:  r.InitialCountryID.Hex(),
		NeighborIDs:       neighbors,
	}
}

func toParty(p *trackers.Party) *model.Party {
	if p == nil {
		return nil
	}
	members := make([]string, len(p.MemberUserIDs))
	for i, id := range p.MemberUserIDs {
		members[i] = id.Hex()
	}
	return &model.Party{
		ID:          p.ID.Hex(),
		Name:        p.Name,
		Description: p.Description,
		AvatarURL:   p.AvatarUrl,
		Ethics: &model.Ethics{
			Unethical:     p.Ethics.Unethical,
			Militarism:    int32(p.Ethics.Militarism),
			Isolationism:  int32(p.Ethics.Isolationism),
			Imperialism:   int32(p.Ethics.Imperialism),
			Industrialism: int32(p.Ethics.Industrialism),
		},
		LastUpdated:   p.LastUpdated,
		CountryID:     p.CountryID.Hex(),
		RegionID:      p.RegionID.Hex(),
		LeaderUserID:  p.LeaderUserID.Hex(),
		MemberUserIDs: members,
	}
}

func toMu(m *trackers.Mu) *model.Mu {
	if m == nil {
		return nil
	}
	members := make([]string, len(m.MemberUserIDs))
	for i, id := range m.MemberUserIDs {
		members[i] = id.Hex()
	}
	return &model.Mu{
		ID:            m.ID.Hex(),
		Name:          m.Name,
		AvatarURL:     m.AvatarUrl,
		Level:         int32(m.Level),
		Hq:            int32(m.HeadQuarterLevel),
		Dorms:         int32(m.DormitoriesLevel),
		MercRep:       m.MercenaryReputation,
		OwnerUserID:   m.OwnerUserID.Hex(),
		RegionID:      m.RegionID.Hex(),
		MemberUserIDs: members,
	}
}

func toAlliance(a *trackers.Alliance) *model.Alliance {
	if a == nil {
		return nil
	}
	return &model.Alliance{
		ID:   a.ID.Hex(),
		Name: a.Name,
	}
}

func toBattle(b *trackers.Battle) *model.Battle {
	if b == nil {
		return nil
	}
	return &model.Battle{
		ID:                 b.ID.Hex(),
		AttackerDamages:    int32(b.AttackerDamages),
		DefenderDamages:    int32(b.DefenderDamages),
		WinnerSide:         (*enums.Side)(b.WinnerSide),
		IsActive:           b.IsActive,
		EndedAt:            b.EndedAt,
		LastUpdated:        b.LastUpdated,
		AttackerCountryID:  b.AttackerCountryID.Hex(),
		DefenderCountryID:  b.DefenderCountryID.Hex(),
		AttackerAllianceID: hexPtr(b.AttackerAllianceID),
		DefenderAllianceID: hexPtr(b.DefenderAllianceID),
		AttackerRegionID:   hexPtr(b.AttackerRegionID),
		DefenderRegionID:   b.DefenderRegionID.Hex(),
	}
}

func toItem(i *trackers.Item) *model.Item {
	if i == nil {
		return nil
	}
	return &model.Item{
		ID:          i.ID.Hex(),
		ItemCode:    i.ItemCode,
		State:       int32(i.State),
		Status:      i.Status,
		OwnerUserID: i.OwnerUserID.Hex(),
		SkillsMap:   i.Skills,
	}
}

func toDamage(d trackers.Damage) *model.Damage {
	return &model.Damage{
		ID:           d.ID.Hex(),
		Side:         d.Side,
		MilitaryRank: int32(d.MilitaryRank),
		Damages:      int32(d.Damages),
		At:           d.At,
		Ammo:         d.Ammo,
		BattleID:     d.BattleID.Hex(),
		UserID:       d.UserID.Hex(),
		CountryID:    d.CountryID.Hex(),
		AllianceID:   hexPtr(d.AllianceID),
		MuID:         hexPtr(d.MuID),
		PartyID:      hexPtr(d.PartyID),
		SkillID:      d.SkillID.Hex(),
		WeaponID:     hexPtr(d.WeaponID),
		HelmetID:     hexPtr(d.HelmetID),
		ChestID:      hexPtr(d.ChestID),
		PantsID:      hexPtr(d.PantsID),
		BootsID:      hexPtr(d.BootsID),
		GlovesID:     hexPtr(d.GlovesID),
	}
}

func toSkill(s *trackers.Skill) *model.Skill {
	if s == nil {
		return nil
	}
	return &model.Skill{
		ID:     s.ID.Hex(),
		Since:  s.Since,
		Set:    toSkillSet(s.Skills),
		UserID: s.UserID.Hex(),
	}
}

func toSkillSet(s trackers.UserSkills) *model.SkillSet {
	return &model.SkillSet{
		Energy:           int32(s.Energy),
		Health:           int32(s.Health),
		Hunger:           int32(s.Hunger),
		Attack:           int32(s.Attack),
		Companies:        int32(s.Companies),
		Entrepreneurship: int32(s.Entrepreneurship),
		Production:       int32(s.Production),
		CriticalChance:   int32(s.CriticalChance),
		CriticalDamages:  int32(s.CriticalDamages),
		Armor:            int32(s.Armor),
		Precision:        int32(s.Precision),
		Dodge:            int32(s.Dodge),
		LootChance:       int32(s.LootChance),
		Management:       int32(s.Management),
	}
}

func toCompany(c *trackers.Company) *model.Company {
	if c == nil {
		return nil
	}
	return &model.Company{
		ID:       c.ID.Hex(),
		Name:     c.Name,
		ItemCode: c.ItemCode,
		UserID:   c.UserID.Hex(),
		RegionID: c.RegionID.Hex(),
	}
}

func toEmployee(e *trackers.Employee) *model.Employee {
	if e == nil {
		return nil
	}
	return &model.Employee{
		ID:                     e.ID.Hex(),
		Wage:                   e.Wage,
		Fidelity:               int32(e.Fidelity),
		JoinedAt:               e.JoinedAt,
		LastFidelityIncreaseAt: e.LastFidelityIncreaseAt,
		UserID:                 e.UserID.Hex(),
		CompanyID:              e.CompanyID.Hex(),
		EmployerID:             e.EmployerID.Hex(),
	}
}

func toTradeOffer(o *trackers.TradeOffer) *model.TradeOffer {
	if o == nil {
		return nil
	}
	return &model.TradeOffer{
		ID:        o.ID.Hex(),
		ItemCode:  o.ItemCode,
		Side:      o.Side,
		Quantity:  int32(o.Quantity),
		Fulfilled: int32(o.Fulfilled),
		Cancelled: o.Cancelled,
		Price:     o.Price,
		Since:     o.Since,
		UserID:    o.UserID.Hex(),
		CountryID: hexPtr(o.CountryID),
		MuID:      hexPtr(o.MuID),
	}
}

// ---------- transaction mappers ----------

func toMarket(t transactions.MarketTransaction) *model.MarketTransaction {
	return &model.MarketTransaction{
		ID:       t.ID.Hex(),
		Money:    t.Money,
		SellerID: t.SellerID.Hex(),
		BuyerID:  t.BuyerID.Hex(),
		ItemID:   t.ItemID.Hex(),
	}
}

func toWage(t transactions.WageTransaction) *model.WageTransaction {
	return &model.WageTransaction{
		ID:               t.ID.Hex(),
		Money:            t.Money,
		ProductionPoints: int32(t.ProductionPoints),
		EmployeeID:       t.EmployeeID.Hex(),
		EmployerID:       t.EmployerID.Hex(),
	}
}

func toCase(t transactions.CaseTransaction) *model.CaseTransaction {
	return &model.CaseTransaction{
		ID:     t.ID.Hex(),
		Case:   t.Case,
		UserID: t.UserID.Hex(),
		ItemID: t.ItemID.Hex(),
	}
}

func toTrade(t transactions.TradeTransaction) *model.TradeTransaction {
	return &model.TradeTransaction{
		ID:              t.ID.Hex(),
		ItemCode:        t.ItemCode,
		Money:           t.Money,
		Quantity:        int32(t.Quantity),
		TimeTillSale:    t.TimeTillSale,
		SellerID:        t.SellerID.Hex(),
		BuyerID:         t.BuyerID.Hex(),
		SellerMuID:      hexPtr(t.SellerMuID),
		BuyerMuID:       hexPtr(t.BuyerMuID),
		SellerCountryID: hexPtr(t.SellerCountryID),
		BuyerCountryID:  hexPtr(t.BuyerCountryID),
	}
}

func toCraft(t transactions.CraftTransaction) *model.CraftTransaction {
	return &model.CraftTransaction{
		ID:         t.ID.Hex(),
		ScrapsCost: int32(t.ScrapsCost),
		UserID:     t.UserID.Hex(),
		ItemID:     t.ItemID.Hex(),
	}
}

func toDismantle(t transactions.DismantleTransaction) *model.DismantleTransaction {
	return &model.DismantleTransaction{
		ID:             t.ID.Hex(),
		ScrapsReceived: int32(t.ScrapsReceived),
		UserID:         t.UserID.Hex(),
		ItemID:         t.ItemID.Hex(),
	}
}

func toLoot(t transactions.LootTransaction) *model.LootTransaction {
	return &model.LootTransaction{
		ID:     t.ID.Hex(),
		UserID: t.UserID.Hex(),
		ItemID: t.ItemID.Hex(),
	}
}

// ---------- event mappers ----------

func toUserNameChange(e events.UserNameChange) *model.UserNameChange {
	return &model.UserNameChange{ID: e.ID.Hex(), At: atOf(e.ID), Username: e.Username, UserID: e.UserID.Hex()}
}

func toUserCountryChange(e events.UserCountryChange) *model.UserCountryChange {
	return &model.UserCountryChange{ID: e.ID.Hex(), At: atOf(e.ID), UserID: e.UserID.Hex(), CountryID: hexPtr(e.CountryID)}
}

func toUserPartyChange(e events.UserPartyChange) *model.UserPartyChange {
	return &model.UserPartyChange{ID: e.ID.Hex(), At: atOf(e.ID), UserID: e.UserID.Hex(), PartyID: hexPtr(e.PartyID)}
}

func toUserCompanyChange(e events.UserCompanyChange) *model.UserCompanyChange {
	return &model.UserCompanyChange{ID: e.ID.Hex(), At: atOf(e.ID), UserID: e.UserID.Hex(), CompanyID: hexPtr(e.CompanyID)}
}

func toUserMuChange(e events.UserMUChange) *model.UserMuChange {
	return &model.UserMuChange{ID: e.ID.Hex(), At: atOf(e.ID), UserID: e.UserID.Hex(), MuID: hexPtr(e.MUID)}
}

func toUserSkillChange(e events.UserSkillChange) *model.UserSkillChange {
	return &model.UserSkillChange{ID: e.ID.Hex(), At: atOf(e.ID), UserID: e.UserID.Hex(), SkillsMap: e.Skills}
}

func toPartyNameChange(e events.PartyNameChange) *model.PartyNameChange {
	return &model.PartyNameChange{ID: e.ID.Hex(), At: atOf(e.ID), Name: e.Name, PartyID: e.PartyID.Hex()}
}

func toPartyLeaderChange(e events.PartyLeaderChange) *model.PartyLeaderChange {
	return &model.PartyLeaderChange{ID: e.ID.Hex(), At: atOf(e.ID), PartyID: e.PartyID.Hex(), LeaderUserID: e.LeaderUserID.Hex()}
}

func toPartyDescriptionChange(e events.PartyDescriptionChange) *model.PartyDescriptionChange {
	return &model.PartyDescriptionChange{ID: e.ID.Hex(), At: atOf(e.ID), Description: e.Description, PartyID: e.PartyID.Hex()}
}

func toPartyEthicsChange(e events.PartyEthicsChange) *model.PartyEthicsChange {
	return &model.PartyEthicsChange{
		ID: e.ID.Hex(), At: atOf(e.ID), PartyID: e.PartyID.Hex(),
		Ethics: &model.Ethics{
			Unethical:     e.Ethics.Unethical,
			Militarism:    int32(e.Ethics.Militarism),
			Isolationism:  int32(e.Ethics.Isolationism),
			Imperialism:   int32(e.Ethics.Imperialism),
			Industrialism: int32(e.Ethics.Industrialism),
		},
	}
}

func toCountryRulingPartyChange(e events.CountryRulingPartyChange) *model.CountryRulingPartyChange {
	return &model.CountryRulingPartyChange{ID: e.ID.Hex(), At: atOf(e.ID), CountryID: e.CountryID.Hex(), PartyID: hexPtr(e.PartyID)}
}

func toCountrySpecialisationChange(e events.CountrySpecialisationChange) *model.CountrySpecialisationChange {
	return &model.CountrySpecialisationChange{ID: e.ID.Hex(), At: atOf(e.ID), ItemCode: e.SpecialisationItemCode, CountryID: e.CountryID.Hex()}
}

func toCountryAllianceJoin(e events.CountryAllianceJoin) *model.CountryAllianceJoin {
	return &model.CountryAllianceJoin{ID: e.ID.Hex(), At: atOf(e.ID), CountryID: e.CountryID.Hex(), AllianceID: e.AllianceID.Hex()}
}

func toCountryAllianceLeave(e events.CountryAllianceLeave) *model.CountryAllianceLeave {
	return &model.CountryAllianceLeave{ID: e.ID.Hex(), At: atOf(e.ID), CountryID: e.CountryID.Hex(), PrevAllianceID: hexPtr(e.PrevAllianceID)}
}

func toRegionOwnerChange(e events.RegionOwnerChange) *model.RegionOwnerChange {
	return &model.RegionOwnerChange{ID: e.ID.Hex(), At: atOf(e.ID), RegionID: e.RegionID.Hex(), CountryID: e.CountryID.Hex()}
}

func toRegionDeposit(e events.RegionDeposit) *model.RegionDeposit {
	return &model.RegionDeposit{
		ID: e.ID.Hex(), At: atOf(e.ID), Type: e.Type, StartsAt: e.StartsAt,
		EndsAt: e.EndsAt, BonusPercent: e.BonusPercent, RegionID: e.RegionID.Hex(),
	}
}

func toRegionStrategicResource(e events.RegionStrategicResource) *model.RegionStrategicResource {
	return &model.RegionStrategicResource{ID: e.ID.Hex(), At: atOf(e.ID), Resource: e.Resource, RegionID: e.RegionID.Hex()}
}

func toCompanyRegionChange(e events.CompanyRegionChange) *model.CompanyRegionChange {
	return &model.CompanyRegionChange{ID: e.ID.Hex(), At: atOf(e.ID), CompanyID: e.CompanyID.Hex(), RegionID: e.RegionID.Hex()}
}

func toCompanyItemCodeChange(e events.CompanyItemCodeChange) *model.CompanyItemCodeChange {
	return &model.CompanyItemCodeChange{ID: e.ID.Hex(), At: atOf(e.ID), ItemCode: e.ItemCode, CompanyID: e.CompanyID.Hex()}
}

func toEmployeeWageChange(e events.EmployeeWageChange) *model.EmployeeWageChange {
	return &model.EmployeeWageChange{ID: e.ID.Hex(), At: atOf(e.ID), Wage: e.Wage, UserID: e.UserID.Hex(), CompanyID: e.CompanyID.Hex()}
}

func toMuNameChange(e events.MuNameChange) *model.MuNameChange {
	return &model.MuNameChange{ID: e.ID.Hex(), At: atOf(e.ID), Name: e.Name, MuID: e.MuID.Hex()}
}

func toMuOwnerChange(e events.MuOwnerChange) *model.MuOwnerChange {
	return &model.MuOwnerChange{ID: e.ID.Hex(), At: atOf(e.ID), MuID: e.MuID.Hex(), OwnerUserID: e.OwnerUserID.Hex()}
}

func toMuMercenaryReputationChange(e events.MuMercenaryReputationChange) *model.MuMercenaryReputationChange {
	return &model.MuMercenaryReputationChange{ID: e.ID.Hex(), At: atOf(e.ID), MercRep: e.MercenaryReputation, MuID: e.MuID.Hex()}
}

func toBattleOrderChange(e events.BattleOrderChange) *model.BattleOrderChange {
	return &model.BattleOrderChange{
		ID: e.ID.Hex(), At: e.At, Side: e.Side, Kind: e.Kind, Action: e.Action,
		BattleID: e.BattleID.Hex(), EntityID: e.EntityID.Hex(),
	}
}

// ---------- processed mappers ----------

func toUserFinanceReport(r processedreports.UserFinanceReport) *model.UserFinanceReport {
	return &model.UserFinanceReport{
		ID: r.ID, DayStart: r.DayStart, WagesPaid: r.WagesPaid, WagesEarned: r.WagesEarned,
		ItemsBought: r.ItemsBought, ItemsSold: r.ItemsSold, EquipBought: r.EquipBought, EquipSold: r.EquipSold,
		ValueDismantled: r.ValueDismantled, CasesOpened: int32(r.CasesOpened), CasesNet: r.CasesNet, UserID: r.UserID.Hex(),
	}
}

func toUserFlipEvent(e processedestimators.UserFlipEvent) *model.UserFlipEvent {
	return &model.UserFlipEvent{
		ID: e.ID.Hex(), ItemCode: e.ItemCode, Quantity: int32(e.Quantity), BuyCost: e.BuyCost,
		SellRevenue: e.SellRevenue, Profit: e.Profit, HeldMs: e.HeldMs, At: e.At, UserID: e.UserID.Hex(),
	}
}

func toUserFlipState(s *processedestimators.UserFlipState) *model.UserFlipState {
	if s == nil {
		return nil
	}
	lots := make(map[string][]model.FlipLot, len(s.OpenLots))
	for code, group := range s.OpenLots {
		conv := make([]model.FlipLot, len(group))
		for i, l := range group {
			conv[i] = model.FlipLot{TradeID: l.TradeID.Hex(), Quantity: int32(l.Quantity), UnitPrice: l.UnitPrice, BoughtAt: l.BoughtAt}
		}
		lots[code] = conv
	}
	return &model.UserFlipState{
		TotalFlips: int32(s.TotalFlips), TotalProfit: s.TotalProfit, UpdatedAt: s.UpdatedAt,
		UserID: s.UserID.Hex(), OpenLotsMap: lots,
	}
}

func toUserInventory(inv *processedestimators.UserInventory) *model.UserInventory {
	if inv == nil {
		return nil
	}
	return &model.UserInventory{UpdatedAt: inv.UpdatedAt, UserID: inv.UserID.Hex(), ItemsMap: inv.Items}
}

func toUserBattleParticipation(p *processedestimators.UserBattleParticipation) *model.UserBattleParticipation {
	if p == nil {
		return nil
	}
	return &model.UserBattleParticipation{
		TotalDamage: p.TotalDamage, BattlesParticipated: int32(p.BattlesParticipated), NegativeDamage: p.NegativeDamage,
		OwnCountryBattles: int32(p.OwnCountryBattles), OwnCountryParticipated: int32(p.OwnCountryParticipated),
		MuOrderBattles: int32(p.MuOrderBattles), MuOrderParticipated: int32(p.MuOrderParticipated),
		UpdatedAt: p.UpdatedAt, UserID: p.UserID.Hex(),
	}
}

func toCountryTaxFlow(f processedreports.CountryTaxFlow) *model.CountryTaxFlow {
	hijackers := make([]*model.TaxHijack, len(f.Hijackers))
	for i, h := range f.Hijackers {
		hijackers[i] = &model.TaxHijack{Amount: h.Amount, CountryID: h.CountryID.Hex()}
	}
	sources := make([]*model.TaxSource, len(f.Sources))
	for i, s := range f.Sources {
		sources[i] = &model.TaxSource{Total: s.Total, Hijacked: s.Hijacked, CorePct: s.CorePct, CountryID: s.CountryID.Hex()}
	}
	return &model.CountryTaxFlow{
		ID: f.ID, HourStart: f.HourStart, TotalTax: f.TotalTax, HijackedIn: f.HijackedIn,
		CoreEarned: f.CoreEarned, NonCoreEarned: f.NonCoreEarned, HijackedOut: f.HijackedOut,
		Hijackers: hijackers, Sources: sources, CountryID: f.CountryID.Hex(),
	}
}

func toCountryMoneyFlowReport(r processedreports.CountryMoneyFlowReport) *model.CountryMoneyFlowReport {
	counterparts := make([]*model.CountryMoneyFlowCounterpart, len(r.Counterparts))
	for i, cp := range r.Counterparts {
		counterparts[i] = &model.CountryMoneyFlowCounterpart{
			InEquipment:  cp.InEquipment,
			OutEquipment: cp.OutEquipment,
			InItems:      cp.InItems,
			OutItems:     cp.OutItems,
			InWages:      cp.InWages,
			OutWages:     cp.OutWages,
			CountryID:    cp.CountryID.Hex(),
		}
	}
	return &model.CountryMoneyFlowReport{
		ID:                      r.ID,
		DayStart:                r.DayStart,
		InEquipment:             r.InEquipment,
		OutEquipment:            r.OutEquipment,
		InItems:                 r.InItems,
		OutItems:                r.OutItems,
		InWages:                 r.InWages,
		OutWages:                r.OutWages,
		InEquipmentDomestic:     r.InEquipmentDomestic,
		OutEquipmentDomestic:    r.OutEquipmentDomestic,
		InItemsDomestic:         r.InItemsDomestic,
		OutItemsDomestic:        r.OutItemsDomestic,
		InWagesDomestic:         r.InWagesDomestic,
		OutWagesDomestic:        r.OutWagesDomestic,
		InEquipmentCrossBorder:  r.InEquipmentCrossBorder,
		OutEquipmentCrossBorder: r.OutEquipmentCrossBorder,
		InItemsCrossBorder:      r.InItemsCrossBorder,
		OutItemsCrossBorder:     r.OutItemsCrossBorder,
		InWagesCrossBorder:      r.InWagesCrossBorder,
		OutWagesCrossBorder:     r.OutWagesCrossBorder,
		Counterparts:            counterparts,
		CountryID:               r.CountryID.Hex(),
	}
}

func toMuCountryMoneyFlowReport(r processedreports.MuCountryMoneyFlowReport) *model.MuCountryMoneyFlowReport {
	counterparts := make([]*model.MuCountryMoneyFlowCounterpart, len(r.Counterparts))
	for i, cp := range r.Counterparts {
		counterparts[i] = &model.MuCountryMoneyFlowCounterpart{
			InEquipment:  cp.InEquipment,
			OutEquipment: cp.OutEquipment,
			InItems:      cp.InItems,
			OutItems:     cp.OutItems,
			InWages:      cp.InWages,
			OutWages:     cp.OutWages,
			CountryID:    cp.CountryID.Hex(),
		}
	}
	return &model.MuCountryMoneyFlowReport{
		ID:                                      r.ID,
		DayStart:                                r.DayStart,
		InEquipment:                             r.InEquipment,
		OutEquipment:                            r.OutEquipment,
		InItems:                                 r.InItems,
		OutItems:                                r.OutItems,
		InWages:                                 r.InWages,
		OutWages:                                r.OutWages,
		InEquipmentInsideMu:                     r.InEquipmentInsideMu,
		OutEquipmentInsideMu:                    r.OutEquipmentInsideMu,
		InItemsInsideMu:                         r.InItemsInsideMu,
		OutItemsInsideMu:                        r.OutItemsInsideMu,
		InWagesInsideMu:                         r.InWagesInsideMu,
		OutWagesInsideMu:                        r.OutWagesInsideMu,
		InEquipmentSameCountryOutsideMu:         r.InEquipmentSameCountryOutsideMu,
		OutEquipmentSameCountryOutsideMu:        r.OutEquipmentSameCountryOutsideMu,
		InItemsSameCountryOutsideMu:             r.InItemsSameCountryOutsideMu,
		OutItemsSameCountryOutsideMu:            r.OutItemsSameCountryOutsideMu,
		InWagesSameCountryOutsideMu:             r.InWagesSameCountryOutsideMu,
		OutWagesSameCountryOutsideMu:            r.OutWagesSameCountryOutsideMu,
		InEquipmentCrossBorderOutsideMuCountry:  r.InEquipmentCrossBorderOutsideMuCountry,
		OutEquipmentCrossBorderOutsideMuCountry: r.OutEquipmentCrossBorderOutsideMuCountry,
		InItemsCrossBorderOutsideMuCountry:      r.InItemsCrossBorderOutsideMuCountry,
		OutItemsCrossBorderOutsideMuCountry:     r.OutItemsCrossBorderOutsideMuCountry,
		InWagesCrossBorderOutsideMuCountry:      r.InWagesCrossBorderOutsideMuCountry,
		OutWagesCrossBorderOutsideMuCountry:     r.OutWagesCrossBorderOutsideMuCountry,
		Counterparts:                            counterparts,
		MuID:                                    r.MuID.Hex(),
	}
}

func toCountryFlipEvent(e processedestimators.CountryFlipEvent) *model.CountryFlipEvent {
	return &model.CountryFlipEvent{
		ID: e.ID.Hex(), ItemCode: e.ItemCode, Quantity: int32(e.Quantity), BuyCost: e.BuyCost,
		SellRevenue: e.SellRevenue, Profit: e.Profit, At: e.At, CountryID: e.CountryID.Hex(),
	}
}

func toCountryFlipState(s *processedestimators.CountryFlipState) *model.CountryFlipState {
	if s == nil {
		return nil
	}
	return &model.CountryFlipState{
		TotalTrades: int32(s.TotalTrades), Profitable: int32(s.Profitable),
		TotalProfit: s.TotalProfit, UpdatedAt: s.UpdatedAt, CountryID: s.CountryID.Hex(),
	}
}

func toCountryInventory(inv *processedestimators.CountryInventory) *model.CountryInventory {
	if inv == nil {
		return nil
	}
	lots := make(map[string][]model.InventoryLot, len(inv.Lots))
	for code, group := range inv.Lots {
		conv := make([]model.InventoryLot, len(group))
		for i, l := range group {
			conv[i] = model.InventoryLot{TradeID: l.TradeID.Hex(), Quantity: int32(l.Quantity), UnitPrice: l.UnitPrice, BoughtAt: l.BoughtAt}
		}
		lots[code] = conv
	}
	return &model.CountryInventory{UpdatedAt: inv.UpdatedAt, CountryID: inv.CountryID.Hex(), LotsMap: lots}
}

func toEntityWealthReport(r processedreports.EntityWealthReport) *model.EntityWealthReport {
	return &model.EntityWealthReport{
		ID: r.ID, DayStart: r.DayStart, MemberCount: int32(r.MemberCount), TotalDamage: r.TotalDamage,
		TotalWealth: r.TotalWealth, WagesPaid: r.WagesPaid, WagesEarned: r.WagesEarned,
		EntityType: r.EntityType, EntityID: r.EntityID.Hex(),
	}
}

func toBattleDamageReport(r processedreports.BattleDamageReport) *model.BattleDamageReport {
	equip := make([]*model.EquipmentUsage, len(r.Equipment))
	for i, e := range r.Equipment {
		equip[i] = &model.EquipmentUsage{ItemCode: e.ItemCode, Count: e.Count, Value: e.Value}
	}
	return &model.BattleDamageReport{
		ID: r.ID, IntervalStart: r.IntervalStart, Side: r.Side, Damage: r.Damage, DamagePct: r.DamagePct,
		Equipment: equip, BattleID: r.BattleID.Hex(), EntityType: r.EntityType, EntityID: r.EntityID.Hex(),
	}
}

func toItemCandle(c processedcandles.ItemCandle) *model.ItemCandle {
	return &model.ItemCandle{
		ID: c.ID, ItemCode: c.ItemCode, BucketStart: c.BucketStart, Open: c.Open, High: c.High,
		Low: c.Low, Close: c.Close, Avg: c.Avg, Volume: int32(c.Volume), Money: c.Money, Count: int32(c.Count),
	}
}

func toWageCandle(c processedcandles.WageCandle) *model.WageCandle {
	return &model.WageCandle{
		ID: c.ID, BucketStart: c.BucketStart, Open: c.Open, High: c.High, Low: c.Low, Close: c.Close,
		Avg: c.Avg, Volume: int32(c.Volume), Money: c.Money, Count: int32(c.Count),
	}
}

func toMarketState(s processedreports.MarketState) *model.MarketState {
	return &model.MarketState{
		At: s.At, AvgWage24h: s.AvgWage24h, WageVolume24h: int32(s.WageVolume24h), MarketVolume24h: s.MarketVolume24h,
		WageMin: s.WageMin, WageMax: s.WageMax, WageAvgWeighted: s.WageAvgWeighted,
	}
}

func toWageMarketState(s processedreports.WageMarketState) *model.WageMarketState {
	least := make([]*model.WagePaidUser, len(s.Top10Least24h))
	for i, u := range s.Top10Least24h {
		least[i] = &model.WagePaidUser{TotalPaid: u.TotalPaid, UserID: u.UserID.Hex()}
	}
	most := make([]*model.WagePaidUser, len(s.Top10Most24h))
	for i, u := range s.Top10Most24h {
		most[i] = &model.WagePaidUser{TotalPaid: u.TotalPaid, UserID: u.UserID.Hex()}
	}
	return &model.WageMarketState{
		At: s.At, AvgWeighted14d: s.AvgWeighted14d, Min14d: s.Min14d, Max14d: s.Max14d,
		TotalPaid14d: s.TotalPaid14d, Top10Least24h: least, Top10Most24h: most,
	}
}

func toInflationPoint(p processedestimators.InflationPoint) *model.InflationPoint {
	return &model.InflationPoint{
		ID: p.ID, DayStart: p.DayStart, IndexValue: p.IndexValue, PctChange: p.PctChange,
		PerItem: floatEntries(p.PerItem),
	}
}

func toDismantleReport(r processedreports.DismantleReport) *model.DismantleReport {
	return &model.DismantleReport{
		ID: r.ID, HourStart: r.HourStart, Count: int32(r.Count), StateBuckets: intEntries(r.StateBuckets),
	}
}

func toItemMarketReport(r *processedreports.ItemMarketReport) *model.ItemMarketReport {
	if r == nil {
		return nil
	}
	bids := make([]*model.OrderbookLevel, len(r.Bids))
	for i, l := range r.Bids {
		bids[i] = &model.OrderbookLevel{Price: l.Price, Quantity: int32(l.Quantity)}
	}
	asks := make([]*model.OrderbookLevel, len(r.Asks))
	for i, l := range r.Asks {
		asks[i] = &model.OrderbookLevel{Price: l.Price, Quantity: int32(l.Quantity)}
	}
	buy := make([]*model.EffectivePrice, len(r.EffectiveBuy))
	for i, p := range r.EffectiveBuy {
		buy[i] = &model.EffectivePrice{Size: int32(p.Size), AvgPrice: p.AvgPrice}
	}
	sell := make([]*model.EffectivePrice, len(r.EffectiveSell))
	for i, p := range r.EffectiveSell {
		sell[i] = &model.EffectivePrice{Size: int32(p.Size), AvgPrice: p.AvgPrice}
	}
	return &model.ItemMarketReport{
		ItemCode: r.ItemCode, Volume24h: int32(r.Volume24h), AvgWeighted24h: r.AvgWeighted24h,
		PctChange24h: r.PctChange24h, Low24h: r.Low24h, High24h: r.High24h, Spread: r.Spread,
		Bids: bids, Asks: asks, EffectiveBuy: buy, EffectiveSell: sell, UpdatedAt: r.UpdatedAt,
	}
}

func toCasesReport(r processedreports.CasesReport) *model.CasesReport {
	items := make([]*model.CaseItemStat, len(r.PerItem))
	for i, s := range r.PerItem {
		items[i] = &model.CaseItemStat{
			ItemCode: s.ItemCode, AvgWeighted14d: s.AvgWeighted14d, TotalDrops: int32(s.TotalDrops),
			PerSkillRoll: intEntries(s.PerSkillRoll),
		}
	}
	return &model.CasesReport{
		Case: r.Case, TotalOpened: int32(r.TotalOpened), UniqueItemCodes: r.UniqueItemCodes,
		ExpectedValue14d: r.ExpectedValue14d, PerItem: items, UpdatedAt: r.UpdatedAt,
	}
}

// ---------- arg helpers ----------

// entityKindStr maps an EntityKind enum to the lowercase entityType stored in reports.
func entityKindStr(k *model.EntityKind) *string {
	if k == nil {
		return nil
	}
	var s string
	switch *k {
	case model.EntityKindUser:
		s = "user"
	case model.EntityKindCountry:
		s = "country"
	case model.EntityKindParty:
		s = "party"
	case model.EntityKindMu:
		s = "mu"
	case model.EntityKindAlliance:
		s = "alliance"
	case model.EntityKindSide:
		s = "side"
	default:
		return nil
	}
	return &s
}

func toBattleFilter(f *model.BattleFilter) trackers.BattleFilter {
	if f == nil {
		return trackers.BattleFilterAll
	}
	switch *f {
	case model.BattleFilterActive:
		return trackers.BattleFilterActive
	case model.BattleFilterFinalized:
		return trackers.BattleFilterFinalized
	default:
		return trackers.BattleFilterAll
	}
}

// battleMatchesFilter reports whether a mapped battle satisfies the filter, used
// where battles are reached via a join and filtered after loading.
func battleMatchesFilter(b *model.Battle, f trackers.BattleFilter) bool {
	switch f {
	case trackers.BattleFilterActive:
		return b.IsActive
	case trackers.BattleFilterFinalized:
		return b.WinnerSide != nil
	default:
		return true
	}
}

package graph

import (
	"context"

	"github.com/warerastats/api/graph/model"
)

func (r *battleDamageReportResolver) Battle(ctx context.Context, obj *model.BattleDamageReport) (*model.Battle, error) {
	return loadBattle(ctx, obj.BattleID)
}

func (r *battleDamageReportResolver) Entity(ctx context.Context, obj *model.BattleDamageReport) (model.Entity, error) {
	return loadEntity(ctx, obj.EntityType, obj.EntityID)
}

func (r *countryFlipEventResolver) Country(ctx context.Context, obj *model.CountryFlipEvent) (*model.Country, error) {
	return loadCountry(ctx, obj.CountryID)
}

func (r *countryFlipStateResolver) Country(ctx context.Context, obj *model.CountryFlipState) (*model.Country, error) {
	return loadCountry(ctx, obj.CountryID)
}

func (r *countryInventoryResolver) Lots(ctx context.Context, obj *model.CountryInventory) ([]*model.InventoryLotGroup, error) {
	out := make([]*model.InventoryLotGroup, 0, len(obj.LotsMap))
	for code, lots := range obj.LotsMap {
		group := &model.InventoryLotGroup{ItemCode: code, Lots: make([]*model.InventoryLot, len(lots))}
		for i := range lots {
			l := lots[i]
			group.Lots[i] = &l
		}
		out = append(out, group)
	}
	return out, nil
}

func (r *countryInventoryResolver) Country(ctx context.Context, obj *model.CountryInventory) (*model.Country, error) {
	return loadCountry(ctx, obj.CountryID)
}

func (r *countryTaxFlowResolver) Country(ctx context.Context, obj *model.CountryTaxFlow) (*model.Country, error) {
	return loadCountry(ctx, obj.CountryID)
}

func (r *entityWealthReportResolver) Entity(ctx context.Context, obj *model.EntityWealthReport) (model.Entity, error) {
	return loadEntity(ctx, obj.EntityType, obj.EntityID)
}

func (r *inflationPointResolver) PctChange24h(ctx context.Context, obj *model.InflationPoint) (float64, error) {
	return obj.PctChange, nil
}

func (r *taxHijackResolver) Country(ctx context.Context, obj *model.TaxHijack) (*model.Country, error) {
	return loadCountry(ctx, obj.CountryID)
}

func (r *taxSourceResolver) Country(ctx context.Context, obj *model.TaxSource) (*model.Country, error) {
	return loadCountry(ctx, obj.CountryID)
}

func (r *userBattleParticipationResolver) User(ctx context.Context, obj *model.UserBattleParticipation) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *userFinanceReportResolver) User(ctx context.Context, obj *model.UserFinanceReport) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *userFlipEventResolver) User(ctx context.Context, obj *model.UserFlipEvent) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *userFlipStateResolver) OpenLots(ctx context.Context, obj *model.UserFlipState) ([]*model.FlipLotGroup, error) {
	out := make([]*model.FlipLotGroup, 0, len(obj.OpenLotsMap))
	for code, lots := range obj.OpenLotsMap {
		group := &model.FlipLotGroup{ItemCode: code, Lots: make([]*model.FlipLot, len(lots))}
		for i := range lots {
			l := lots[i]
			group.Lots[i] = &l
		}
		out = append(out, group)
	}
	return out, nil
}

func (r *userFlipStateResolver) User(ctx context.Context, obj *model.UserFlipState) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *userInventoryResolver) Items(ctx context.Context, obj *model.UserInventory) ([]*model.IntEntry, error) {
	return intEntries(obj.ItemsMap), nil
}

func (r *userInventoryResolver) User(ctx context.Context, obj *model.UserInventory) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *wagePaidUserResolver) User(ctx context.Context, obj *model.WagePaidUser) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

// BattleDamageReport returns BattleDamageReportResolver implementation.
func (r *Resolver) BattleDamageReport() BattleDamageReportResolver {
	return &battleDamageReportResolver{r}
}

// CountryFlipEvent returns CountryFlipEventResolver implementation.
func (r *Resolver) CountryFlipEvent() CountryFlipEventResolver { return &countryFlipEventResolver{r} }

// CountryFlipState returns CountryFlipStateResolver implementation.
func (r *Resolver) CountryFlipState() CountryFlipStateResolver { return &countryFlipStateResolver{r} }

// CountryInventory returns CountryInventoryResolver implementation.
func (r *Resolver) CountryInventory() CountryInventoryResolver { return &countryInventoryResolver{r} }

// CountryTaxFlow returns CountryTaxFlowResolver implementation.
func (r *Resolver) CountryTaxFlow() CountryTaxFlowResolver { return &countryTaxFlowResolver{r} }

// EntityWealthReport returns EntityWealthReportResolver implementation.
func (r *Resolver) EntityWealthReport() EntityWealthReportResolver {
	return &entityWealthReportResolver{r}
}

// InflationPoint returns InflationPointResolver implementation.
func (r *Resolver) InflationPoint() InflationPointResolver { return &inflationPointResolver{r} }

// TaxHijack returns TaxHijackResolver implementation.
func (r *Resolver) TaxHijack() TaxHijackResolver { return &taxHijackResolver{r} }

// TaxSource returns TaxSourceResolver implementation.
func (r *Resolver) TaxSource() TaxSourceResolver { return &taxSourceResolver{r} }

// UserBattleParticipation returns UserBattleParticipationResolver implementation.
func (r *Resolver) UserBattleParticipation() UserBattleParticipationResolver {
	return &userBattleParticipationResolver{r}
}

// UserFinanceReport returns UserFinanceReportResolver implementation.
func (r *Resolver) UserFinanceReport() UserFinanceReportResolver {
	return &userFinanceReportResolver{r}
}

// UserFlipEvent returns UserFlipEventResolver implementation.
func (r *Resolver) UserFlipEvent() UserFlipEventResolver { return &userFlipEventResolver{r} }

// UserFlipState returns UserFlipStateResolver implementation.
func (r *Resolver) UserFlipState() UserFlipStateResolver { return &userFlipStateResolver{r} }

// UserInventory returns UserInventoryResolver implementation.
func (r *Resolver) UserInventory() UserInventoryResolver { return &userInventoryResolver{r} }

// WagePaidUser returns WagePaidUserResolver implementation.
func (r *Resolver) WagePaidUser() WagePaidUserResolver { return &wagePaidUserResolver{r} }

type battleDamageReportResolver struct{ *Resolver }
type countryFlipEventResolver struct{ *Resolver }
type countryFlipStateResolver struct{ *Resolver }
type countryInventoryResolver struct{ *Resolver }
type countryTaxFlowResolver struct{ *Resolver }
type entityWealthReportResolver struct{ *Resolver }
type inflationPointResolver struct{ *Resolver }
type taxHijackResolver struct{ *Resolver }
type taxSourceResolver struct{ *Resolver }
type userBattleParticipationResolver struct{ *Resolver }
type userFinanceReportResolver struct{ *Resolver }
type userFlipEventResolver struct{ *Resolver }
type userFlipStateResolver struct{ *Resolver }
type userInventoryResolver struct{ *Resolver }
type wagePaidUserResolver struct{ *Resolver }

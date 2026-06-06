package graph

import (
	"context"

	"github.com/warerastats/api/graph/model"
)

func (r *caseTransactionResolver) User(ctx context.Context, obj *model.CaseTransaction) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *caseTransactionResolver) Item(ctx context.Context, obj *model.CaseTransaction) (*model.Item, error) {
	return loadItemP(ctx, &obj.ItemID)
}

func (r *craftTransactionResolver) User(ctx context.Context, obj *model.CraftTransaction) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *craftTransactionResolver) Item(ctx context.Context, obj *model.CraftTransaction) (*model.Item, error) {
	return loadItemP(ctx, &obj.ItemID)
}

func (r *dismantleTransactionResolver) User(ctx context.Context, obj *model.DismantleTransaction) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *dismantleTransactionResolver) Item(ctx context.Context, obj *model.DismantleTransaction) (*model.Item, error) {
	return loadItemP(ctx, &obj.ItemID)
}

func (r *lootTransactionResolver) User(ctx context.Context, obj *model.LootTransaction) (*model.User, error) {
	return loadUser(ctx, obj.UserID)
}

func (r *lootTransactionResolver) Item(ctx context.Context, obj *model.LootTransaction) (*model.Item, error) {
	return loadItemP(ctx, &obj.ItemID)
}

func (r *marketTransactionResolver) Seller(ctx context.Context, obj *model.MarketTransaction) (*model.User, error) {
	return loadUser(ctx, obj.SellerID)
}

func (r *marketTransactionResolver) Buyer(ctx context.Context, obj *model.MarketTransaction) (*model.User, error) {
	return loadUser(ctx, obj.BuyerID)
}

func (r *marketTransactionResolver) Item(ctx context.Context, obj *model.MarketTransaction) (*model.Item, error) {
	return loadItemP(ctx, &obj.ItemID)
}

func (r *tradeTransactionResolver) Seller(ctx context.Context, obj *model.TradeTransaction) (*model.User, error) {
	return loadUser(ctx, obj.SellerID)
}

func (r *tradeTransactionResolver) Buyer(ctx context.Context, obj *model.TradeTransaction) (*model.User, error) {
	return loadUser(ctx, obj.BuyerID)
}

func (r *tradeTransactionResolver) SellerMu(ctx context.Context, obj *model.TradeTransaction) (*model.Mu, error) {
	return loadMuP(ctx, obj.SellerMuID)
}

func (r *tradeTransactionResolver) BuyerMu(ctx context.Context, obj *model.TradeTransaction) (*model.Mu, error) {
	return loadMuP(ctx, obj.BuyerMuID)
}

func (r *tradeTransactionResolver) SellerCountry(ctx context.Context, obj *model.TradeTransaction) (*model.Country, error) {
	return loadCountryP(ctx, obj.SellerCountryID)
}

func (r *tradeTransactionResolver) BuyerCountry(ctx context.Context, obj *model.TradeTransaction) (*model.Country, error) {
	return loadCountryP(ctx, obj.BuyerCountryID)
}

func (r *wageTransactionResolver) Employee(ctx context.Context, obj *model.WageTransaction) (*model.User, error) {
	return loadUser(ctx, obj.EmployeeID)
}

func (r *wageTransactionResolver) Employer(ctx context.Context, obj *model.WageTransaction) (*model.User, error) {
	return loadUser(ctx, obj.EmployerID)
}

// CaseTransaction returns CaseTransactionResolver implementation.
func (r *Resolver) CaseTransaction() CaseTransactionResolver { return &caseTransactionResolver{r} }

// CraftTransaction returns CraftTransactionResolver implementation.
func (r *Resolver) CraftTransaction() CraftTransactionResolver { return &craftTransactionResolver{r} }

// DismantleTransaction returns DismantleTransactionResolver implementation.
func (r *Resolver) DismantleTransaction() DismantleTransactionResolver {
	return &dismantleTransactionResolver{r}
}

// LootTransaction returns LootTransactionResolver implementation.
func (r *Resolver) LootTransaction() LootTransactionResolver { return &lootTransactionResolver{r} }

// MarketTransaction returns MarketTransactionResolver implementation.
func (r *Resolver) MarketTransaction() MarketTransactionResolver {
	return &marketTransactionResolver{r}
}

// TradeTransaction returns TradeTransactionResolver implementation.
func (r *Resolver) TradeTransaction() TradeTransactionResolver { return &tradeTransactionResolver{r} }

// WageTransaction returns WageTransactionResolver implementation.
func (r *Resolver) WageTransaction() WageTransactionResolver { return &wageTransactionResolver{r} }

type caseTransactionResolver struct{ *Resolver }
type craftTransactionResolver struct{ *Resolver }
type dismantleTransactionResolver struct{ *Resolver }
type lootTransactionResolver struct{ *Resolver }
type marketTransactionResolver struct{ *Resolver }
type tradeTransactionResolver struct{ *Resolver }
type wageTransactionResolver struct{ *Resolver }

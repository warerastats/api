package graph

import (
	"context"

	"github.com/warerastats/api/graph/loaders"
	"github.com/warerastats/api/graph/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// oidOf parses a hex id argument into an ObjectID.
func oidOf(s string) (bson.ObjectID, error) { return bson.ObjectIDFromHex(s) }

// oidArgPtr parses an optional hex id argument into an *ObjectID filter.
func oidArgPtr(s *string) (*bson.ObjectID, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	id, err := bson.ObjectIDFromHex(*s)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// --- batched FK loaders (return nil when the key is absent or empty) ---

func loadUser(ctx context.Context, id string) (*model.User, error) {
	if id == "" {
		return nil, nil
	}
	u, err := loaders.For(ctx).User.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	return toUser(u), nil
}

func loadUserP(ctx context.Context, id *string) (*model.User, error) {
	if id == nil {
		return nil, nil
	}
	return loadUser(ctx, *id)
}

func loadCountry(ctx context.Context, id string) (*model.Country, error) {
	if id == "" {
		return nil, nil
	}
	c, err := loaders.For(ctx).Country.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	return toCountry(c), nil
}

func loadCountryP(ctx context.Context, id *string) (*model.Country, error) {
	if id == nil {
		return nil, nil
	}
	return loadCountry(ctx, *id)
}

func loadParty(ctx context.Context, id string) (*model.Party, error) {
	if id == "" {
		return nil, nil
	}
	p, err := loaders.For(ctx).Party.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	return toParty(p), nil
}

func loadPartyP(ctx context.Context, id *string) (*model.Party, error) {
	if id == nil {
		return nil, nil
	}
	return loadParty(ctx, *id)
}

func loadMu(ctx context.Context, id string) (*model.Mu, error) {
	if id == "" {
		return nil, nil
	}
	m, err := loaders.For(ctx).Mu.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	return toMu(m), nil
}

func loadMuP(ctx context.Context, id *string) (*model.Mu, error) {
	if id == nil {
		return nil, nil
	}
	return loadMu(ctx, *id)
}

func loadRegion(ctx context.Context, id string) (*model.Region, error) {
	if id == "" {
		return nil, nil
	}
	rg, err := loaders.For(ctx).Region.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	return toRegion(rg), nil
}

func loadRegionP(ctx context.Context, id *string) (*model.Region, error) {
	if id == nil {
		return nil, nil
	}
	return loadRegion(ctx, *id)
}

// mapPtr maps a slice of values through a mapper taking *T.
func mapPtr[T, R any](in []T, f func(*T) R) []R {
	out := make([]R, len(in))
	for i := range in {
		out[i] = f(&in[i])
	}
	return out
}

// mapVal maps a slice of values through a mapper taking T.
func mapVal[T, R any](in []T, f func(T) R) []R {
	out := make([]R, len(in))
	for i := range in {
		out[i] = f(in[i])
	}
	return out
}

func loadItem(ctx context.Context, id string) (*model.Item, error) {
	if id == "" {
		return nil, nil
	}
	it, err := loaders.For(ctx).Item.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	return toItem(it), nil
}

func loadItemP(ctx context.Context, id *string) (*model.Item, error) {
	if id == nil || *id == "" {
		return nil, nil
	}
	return loadItem(ctx, *id)
}

func loadCompany(ctx context.Context, id string) (*model.Company, error) {
	if id == "" {
		return nil, nil
	}
	c, err := loaders.For(ctx).Company.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	return toCompany(c), nil
}

func loadCompanyP(ctx context.Context, id *string) (*model.Company, error) {
	if id == nil {
		return nil, nil
	}
	return loadCompany(ctx, *id)
}

func loadBattle(ctx context.Context, id string) (*model.Battle, error) {
	if id == "" {
		return nil, nil
	}
	b, err := loaders.For(ctx).Battle.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	return toBattle(b), nil
}

// loadMembers batch-loads a set of users by hex id, dropping any that are absent.
func (r *Resolver) loadMembers(ctx context.Context, ids []string) ([]*model.User, error) {
	oids := make([]bson.ObjectID, 0, len(ids))
	for _, h := range ids {
		id, err := bson.ObjectIDFromHex(h)
		if err != nil {
			return nil, err
		}
		oids = append(oids, id)
	}
	rows, err := r.Colls.Trackers.User.GetMany(ctx, oids)
	if err != nil {
		return nil, err
	}
	return mapPtr(rows, toUser), nil
}

// loadEntity resolves a polymorphic (entityType, entityId) reference. Returns a
// nil interface when the referenced document is absent so the field stays null.
func loadEntity(ctx context.Context, typ, id string) (model.Entity, error) {
	switch typ {
	case "user":
		u, err := loadUser(ctx, id)
		if err != nil || u == nil {
			return nil, err
		}
		return u, nil
	case "country":
		c, err := loadCountry(ctx, id)
		if err != nil || c == nil {
			return nil, err
		}
		return c, nil
	case "party":
		p, err := loadParty(ctx, id)
		if err != nil || p == nil {
			return nil, err
		}
		return p, nil
	case "mu":
		m, err := loadMu(ctx, id)
		if err != nil || m == nil {
			return nil, err
		}
		return m, nil
	default:
		return nil, nil
	}
}

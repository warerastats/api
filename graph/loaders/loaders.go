// Package loaders provides per-request dataloaders that batch foreign-key
// lookups, collapsing the N+1 queries that naive edge resolution would cause
// (e.g. resolving `.user` on 50 transactions becomes one $in query).
package loaders

import (
	"context"
	"net/http"
	"time"

	"github.com/vikstrous/dataloadgen"
	"github.com/warerastats/models/models"
	"github.com/warerastats/models/models/stores/trackers"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ctxKey struct{}

// Loaders holds one batched loader per FK-resolved entity, scoped to a request.
type Loaders struct {
	User     *dataloadgen.Loader[string, *trackers.User]
	Country  *dataloadgen.Loader[string, *trackers.Country]
	Party    *dataloadgen.Loader[string, *trackers.Party]
	Mu       *dataloadgen.Loader[string, *trackers.Mu]
	Region   *dataloadgen.Loader[string, *trackers.Region]
	Item     *dataloadgen.Loader[string, *trackers.Item]
	Battle   *dataloadgen.Loader[string, *trackers.Battle]
	Company  *dataloadgen.Loader[string, *trackers.Company]
	Alliance *dataloadgen.Loader[string, *trackers.Alliance]
}

// newLoaders builds a fresh loader set bound to the given collections.
func newLoaders(colls *models.Collections) *Loaders {
	opt := dataloadgen.WithWait(time.Millisecond)

	return &Loaders{
		User: dataloadgen.NewLoader(func(ctx context.Context, keys []string) ([]*trackers.User, []error) {
			return batch(ctx, keys, colls.Trackers.User.GetMany, func(u trackers.User) bson.ObjectID { return u.ID })
		}, opt),

		Country: dataloadgen.NewLoader(func(ctx context.Context, keys []string) ([]*trackers.Country, []error) {
			return batch(ctx, keys, getCountries(colls), func(c trackers.Country) bson.ObjectID { return c.ID })
		}, opt),

		Party: dataloadgen.NewLoader(func(ctx context.Context, keys []string) ([]*trackers.Party, []error) {
			return batch(ctx, keys, colls.Trackers.Party.GetMany, func(p trackers.Party) bson.ObjectID { return p.ID })
		}, opt),

		Mu: dataloadgen.NewLoader(func(ctx context.Context, keys []string) ([]*trackers.Mu, []error) {
			return batch(ctx, keys, colls.Trackers.Mu.GetMany, func(m trackers.Mu) bson.ObjectID { return m.ID })
		}, opt),

		Region: dataloadgen.NewLoader(func(ctx context.Context, keys []string) ([]*trackers.Region, []error) {
			return batch(ctx, keys, colls.Trackers.Region.GetMany, func(r trackers.Region) bson.ObjectID { return r.ID })
		}, opt),

		Item: dataloadgen.NewLoader(func(ctx context.Context, keys []string) ([]*trackers.Item, []error) {
			return batch(ctx, keys, colls.Trackers.Item.GetMany, func(i trackers.Item) bson.ObjectID { return i.ID })
		}, opt),

		Battle: dataloadgen.NewLoader(func(ctx context.Context, keys []string) ([]*trackers.Battle, []error) {
			return batch(ctx, keys, colls.Trackers.Battle.GetMany, func(b trackers.Battle) bson.ObjectID { return b.ID })
		}, opt),

		Company: dataloadgen.NewLoader(func(ctx context.Context, keys []string) ([]*trackers.Company, []error) {
			return batch(ctx, keys, colls.Trackers.Company.GetMany, func(c trackers.Company) bson.ObjectID { return c.ID })
		}, opt),

		Alliance: dataloadgen.NewLoader(func(ctx context.Context, keys []string) ([]*trackers.Alliance, []error) {
			return batch(ctx, keys, colls.Trackers.Alliance.GetMany, func(a trackers.Alliance) bson.ObjectID { return a.ID })
		}, opt),
	}
}

// getCountries adapts the Country store, which has no GetMany, into the batch shape.
func getCountries(colls *models.Collections) func(context.Context, []bson.ObjectID) ([]trackers.Country, error) {
	return func(ctx context.Context, ids []bson.ObjectID) ([]trackers.Country, error) {
		all, err := colls.Trackers.Country.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		want := make(map[bson.ObjectID]struct{}, len(ids))
		for _, id := range ids {
			want[id] = struct{}{}
		}
		out := all[:0]
		for _, c := range all {
			if _, ok := want[c.ID]; ok {
				out = append(out, c)
			}
		}
		return out, nil
	}
}

// batch parses hex keys, fetches the matching docs in one query, and returns
// per-key pointers aligned to keys (nil for a missing or invalid key).
func batch[T any](
	ctx context.Context,
	keys []string,
	getMany func(context.Context, []bson.ObjectID) ([]T, error),
	idOf func(T) bson.ObjectID,
) ([]*T, []error) {
	ids := make([]bson.ObjectID, 0, len(keys))
	for _, k := range keys {
		id, err := bson.ObjectIDFromHex(k)
		if err == nil {
			ids = append(ids, id)
		}
	}

	docs, err := getMany(ctx, ids)
	if err != nil {
		errs := make([]error, len(keys))
		for i := range errs {
			errs[i] = err
		}
		return make([]*T, len(keys)), errs
	}

	byID := make(map[bson.ObjectID]*T, len(docs))
	for i := range docs {
		d := docs[i]
		byID[idOf(d)] = &d
	}

	out := make([]*T, len(keys))
	for i, k := range keys {
		id, err := bson.ObjectIDFromHex(k)
		if err != nil {
			continue
		}
		out[i] = byID[id] // nil when absent
	}
	return out, nil
}

// Middleware injects a fresh per-request loader set into the request context.
func Middleware(colls *models.Collections, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxKey{}, newLoaders(colls))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// For returns the loaders bound to the current request.
func For(ctx context.Context) *Loaders {
	return ctx.Value(ctxKey{}).(*Loaders)
}

package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
)

// playgroundKey marks a request as coming from the unauthenticated public
// playground, which activates the tighter time-window limits enforced inside
// the resolvers.
type playgroundKey struct{}

// WithPlayground tags ctx as an unauthenticated playground request so that
// resolvers apply the public time-window caps. Authenticated requests never
// carry this flag and stay unrestricted.
func WithPlayground(ctx context.Context) context.Context {
	return context.WithValue(ctx, playgroundKey{}, true)
}

func isPlayground(ctx context.Context) bool {
	v, _ := ctx.Value(playgroundKey{}).(bool)
	return v
}

// Time-window caps (in days) for the public playground.
const (
	rawTimeWindowDays    = 2  // raw event feeds (flip events)
	reportTimeWindowDays = 14 // aggregated reports / candles / leaderboards
)

// enforceTimeWindow rejects overly wide (from, to) ranges on playground
// requests. Authenticated requests are not restricted.
func enforceTimeWindow(ctx context.Context, from, to time.Time, maxDays int) error {
	if !isPlayground(ctx) {
		return nil
	}
	if to.Before(from) {
		return fmt.Errorf("invalid time window: 'to' must not be before 'from'")
	}
	if to.Sub(from) > time.Duration(maxDays)*24*time.Hour {
		return fmt.Errorf("time window too large: at most %d days are allowed on the public playground", maxDays)
	}
	return nil
}

// PlaygroundGuard returns a gqlgen operation middleware that protects the public
// playground from abusive queries. It rejects queries that exceed maxDepth, pass
// a `first` page size above maxFirst, or whose estimated fan-out (the product of
// page sizes along any path) exceeds maxCost. Introspection queries are exempt so
// the playground can still load the schema.
func PlaygroundGuard(maxDepth, maxFirst, maxCost int) graphql.OperationMiddleware {
	return func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		oc := graphql.GetOperationContext(ctx)
		if oc != nil && oc.Operation != nil {
			cost := 0
			if err := walkSelections(oc, oc.Operation.SelectionSet, 1, maxDepth, maxFirst, 1, &cost); err != nil {
				return errResponse(err.Error())
			}
			if cost > maxCost {
				return errResponse(fmt.Sprintf("query is too expensive for the public playground (estimated cost %d, limit %d); narrow your selection or page sizes", cost, maxCost))
			}
		}
		return next(ctx)
	}
}

// errResponse turns a guard failure into a single-error GraphQL response handler.
func errResponse(msg string) graphql.ResponseHandler {
	return func(ctx context.Context) *graphql.Response {
		return graphql.ErrorResponse(ctx, "%s", msg)
	}
}

// unboundedListSize is the assumed page size for list fields that take no
// `first` argument (e.g. countries, regions, members) when estimating cost.
const unboundedListSize = 100

// walkSelections recursively validates depth, page sizes and accumulates an
// estimated fan-out cost. multiplier is the product of page sizes on the path so
// far; each selected field adds multiplier to the running cost.
func walkSelections(oc *graphql.OperationContext, set ast.SelectionSet, depth, maxDepth, maxFirst, multiplier int, cost *int) error {
	if depth > maxDepth {
		return fmt.Errorf("query is too deeply nested for the public playground (max depth %d)", maxDepth)
	}
	for _, sel := range set {
		switch s := sel.(type) {
		case *ast.Field:
			// Introspection fields are exempt so the playground can load the schema.
			if strings.HasPrefix(s.Name, "__") {
				continue
			}

			size, err := pageSize(oc, s, maxFirst)
			if err != nil {
				return err
			}

			childMult := multiplier
			if isListField(s) {
				childMult *= size
			}
			*cost += childMult

			if len(s.SelectionSet) > 0 {
				if err := walkSelections(oc, s.SelectionSet, depth+1, maxDepth, maxFirst, childMult, cost); err != nil {
					return err
				}
			}
		case *ast.InlineFragment:
			// Fragments do not add a nesting level of their own.
			if err := walkSelections(oc, s.SelectionSet, depth, maxDepth, maxFirst, multiplier, cost); err != nil {
				return err
			}
		case *ast.FragmentSpread:
			if s.Definition != nil {
				if err := walkSelections(oc, s.Definition.SelectionSet, depth, maxDepth, maxFirst, multiplier, cost); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// isListField reports whether the field's schema type is a list.
func isListField(f *ast.Field) bool {
	if f.Definition == nil || f.Definition.Type == nil {
		return false
	}
	return f.Definition.Type.Elem != nil
}

// pageSize resolves the effective page size for a field: its `first` argument
// when present (rejecting values above maxFirst), otherwise a default. Returns
// an error if the requested page size exceeds maxFirst.
func pageSize(oc *graphql.OperationContext, f *ast.Field, maxFirst int) (int, error) {
	for _, a := range f.Arguments {
		if a.Name != "first" {
			continue
		}
		n, ok, err := intArg(oc, a.Value)
		if err != nil {
			return 0, err
		}
		if !ok {
			break
		}
		if n > maxFirst {
			return 0, fmt.Errorf("page size too large on %q (requested %d, max %d on the public playground)", f.Name, n, maxFirst)
		}
		if n <= 0 {
			return unboundedListSize, nil
		}
		return n, nil
	}
	return unboundedListSize, nil
}

// intArg resolves an argument value (literal or variable) to an int.
func intArg(oc *graphql.OperationContext, v *ast.Value) (int, bool, error) {
	if v == nil {
		return 0, false, nil
	}
	if v.Kind == ast.Variable {
		raw, ok := oc.Variables[v.Raw]
		if !ok {
			return 0, false, nil
		}
		return coerceInt(raw)
	}
	if v.Kind == ast.IntValue {
		n, err := strconv.Atoi(v.Raw)
		if err != nil {
			return 0, false, fmt.Errorf("invalid 'first' value %q", v.Raw)
		}
		return n, true, nil
	}
	return 0, false, nil
}

func coerceInt(raw any) (int, bool, error) {
	switch n := raw.(type) {
	case int:
		return n, true, nil
	case int32:
		return int(n), true, nil
	case int64:
		return int(n), true, nil
	case float64:
		return int(n), true, nil
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			return 0, false, fmt.Errorf("invalid 'first' variable %q", n.String())
		}
		return int(i), true, nil
	case string:
		i, err := strconv.Atoi(n)
		if err != nil {
			return 0, false, fmt.Errorf("invalid 'first' variable %q", n)
		}
		return i, true, nil
	default:
		return 0, false, nil
	}
}

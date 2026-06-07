package graph

import (
	"github.com/warerastats/api/graph/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// userCursor resolves a user object id and the optional keyset cursor together.
func (r *userResolver) userCursor(obj *model.User, after *string) (bson.ObjectID, *bson.ObjectID, error) {
	uid, err := oidOf(obj.ID)
	if err != nil {
		return bson.ObjectID{}, nil, err
	}
	before, err := cursorPtr(after)
	if err != nil {
		return bson.ObjectID{}, nil, err
	}
	return uid, before, nil
}

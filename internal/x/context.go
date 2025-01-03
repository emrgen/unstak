package x

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

type SpaceID interface {
	SpaceID() uuid.UUID
}

func GetSpaceID(ctx context.Context) (uuid.UUID, error) {
	spaceID, ok := ctx.Value("space_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("space_id not found")
	}
	return spaceID, nil
}

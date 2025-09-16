package repository

import (
	"context"

	"github.com/SomeHowMicroservice/user/model"
)

type ProfileRepository interface {
	Update(ctx context.Context, id string, updateData map[string]any) error

	FindByID(ctx context.Context, id string) (*model.Profile, error)
}
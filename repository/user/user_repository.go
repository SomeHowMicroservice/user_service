package repository

import (
	"context"

	"github.com/SomeHowMicroservice/user/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error 

	ExistsByUsername(ctx context.Context, username string) (bool, error)

	ExistsById(ctx context.Context, id string) (bool, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)

	FindByUsernameWithProfileAndRoles(ctx context.Context, username string) (*model.User, error)

	FindByEmail(ctx context.Context, email string) (*model.User, error)

	Update(ctx context.Context, id string, updateData map[string]any) error

	FindAllByID(ctx context.Context, ids []string) ([]*model.User, error)

	FindByIDWithProfileAndRoles(ctx context.Context, id string) (*model.User, error)
}
package bot

import (
	"context"

	"github.com/HeadHunter483/go-tg-bot/internal/models"
)

type Repository interface {
	GetUsersCount(context.Context) (int, error)
	GetUserByChatID(context.Context, int64) (models.User, error)
	AddUser(context.Context, models.User) error
}

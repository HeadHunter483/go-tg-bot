package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"

	"github.com/HeadHunter483/go-tg-bot/internal/models"
)

// GetUsersCount gets the amount of users registered in the bot.
func (r *repository) GetUsersCount(
	ctx context.Context,
) (count int, err error) {
	query := `SELECT COUNT(id) FROM users;`
	err = r.pool.QueryRow(ctx, query).Scan(&count)
	return
}

// GetUserByChatID gets user registered in the bot by his chat_id.
func (r *repository) GetUserByChatID(
	ctx context.Context, ChatID int64,
) (user models.User, err error) {
	query := `SELECT * FROM users WHERE chat_id=$1;`
	err = r.pool.QueryRow(ctx, query, ChatID).Scan(
		&user.ID, &user.ChatID, &user.UserName, &user.FirstName,
		&user.LastName, &user.DateRegistered,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		err = ErrNotFound
		return
	}
	return
}

// AddUser adds user to the bot db.
func (r *repository) AddUser(
	ctx context.Context, user models.User,
) (err error) {
	query := `
	INSERT INTO users(
		chat_id, username, first_name, last_name
	) 
	VALUES($1, $2, $3, $4);
	`
	_, err = r.pool.Exec(
		ctx, query,
		user.ChatID,
		user.UserName,
		user.FirstName,
		user.LastName,
	)
	return
}

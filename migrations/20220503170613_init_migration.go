package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upInitMigration, downInitMigration)
}

func upInitMigration(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	const query = `
	CREATE TABLE users(
		ID SERIAL PRIMARY KEY,
		CHAT_ID INT,
		USERNAME TEXT,
		FIRST_NAME TEXT,
		LAST_NAME TEXT,
		DATE_REGISTERED TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := tx.Exec(query)
	return err
}

func downInitMigration(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	const query = `
	DROP TABLE users;
	`
	_, err := tx.Exec(query)
	return err
}

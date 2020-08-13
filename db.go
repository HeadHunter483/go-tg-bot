package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	ID int64
	ChatID int
	UserName string
	FirstName string
	LastName string
	DateRegistered time.Time
}

func getUsersCount() (int, error) {
	// get amount of users registered in the bot
	count := 0
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return count, err
	}
	defer db.Close()

	row := db.QueryRow(`SELECT COUNT(id) FROM users;`)
	err = row.Scan(&count)
	if err != nil {
		return count, err
	}

	return count, err
}

func getUserByChatID(ChatID int64) (*User, error) {
	// get user registered in the bot by his chat_id
	db, err := sql.Open("postgres", dbInfo)
	user := User{}

	if err != nil {
		return &user, err
	}
	defer db.Close()

	row := db.QueryRow(
		`SELECT * FROM users WHERE chat_id=$1;`, ChatID)
	err = row.Scan(&user.ID, &user.ChatID, &user.UserName, &user.FirstName, 
			&user.LastName, &user.DateRegistered)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("User with chat_id=%d not found.", user.ChatID)
	}

	return &user, err
}

func addUser(user *UserData) error {
	// add user to the bot db
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	data := `INSERT INTO users(chat_id, username, first_name, last_name) 
		VALUES($1, $2, $3, $4);`
	if _, err = db.Exec(data, user.ChatID, user.UserName, user.FirstName,
						user.LastName); err != nil {
		return err
	}
	return nil
}

func createTables() error {
	// create tables used in the bot
	log.Print("creating tables")
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	users := `CREATE TABLE users(
		ID SERIAL PRIMARY KEY,
		CHAT_ID INT,
		USERNAME TEXT,
		FIRST_NAME TEXT,
		LAST_NAME TEXT,
		DATE_REGISTERED TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`

	if _, err = db.Exec(users); err != nil {
		return err
	}
	return nil
}

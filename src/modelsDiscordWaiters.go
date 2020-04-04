package main

import (
	"database/sql"
	"time"
)

type DiscordWaiters struct{}

// Waiter is a person who is on the waiting list for the next game
// (they used the "/next" Discord command)
type Waiter struct {
	Username        string
	DiscordMention  string
	DatetimeExpired time.Time
}

func (*DiscordWaiters) GetAll() ([]*Waiter, error) {
	waiters := make([]*Waiter, 0)

	rows, err := db.Query(`
		SELECT
			username,
			discord_mention,
			datetime_expired
		FROM discord_waiters
	`)

	for rows.Next() {
		var waiter Waiter
		if err2 := rows.Scan(
			&waiter.Username,
			&waiter.DiscordMention,
			&waiter.DatetimeExpired,
		); err2 != nil {
			return nil, err2
		}
		waiters = append(waiters, &waiter)
	}

	if err == sql.ErrNoRows {
		return waiters, nil
	}
	if rows.Err() != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	return waiters, nil
}

func (*DiscordWaiters) Insert(waiter *Waiter) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO discord_waiters (username, discord_mention, datetime_expired)
		VALUES (?, ?, ?)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(waiter.Username, waiter.DiscordMention, waiter.DatetimeExpired)
	return err
}

func (*DiscordWaiters) Delete(username string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		DELETE FROM discord_waiters
		WHERE username = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec(username)
	return err
}

func (*DiscordWaiters) DeleteAll() error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		DELETE FROM discord_waiters
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	_, err := stmt.Exec()
	return err
}

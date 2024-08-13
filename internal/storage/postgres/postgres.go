package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	db *pgx.Conn
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.postgres.New"
	ctx := context.Background()

	conn, err := pgx.Connect(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	// TODO add sql injections
	_, err = conn.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS houses(
	    id INTEGER PRIMARY KEY,
	    address TEXT NOT NULL,
		year_build INTEGER NOT NULL CHECK (year_build >= 0),
		developer TEXT NULL,
		created_fl TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_flat_added TIMESTAMP NULL);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = conn.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS users(
	    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	    email TEXT NOT NULL,
	    password TEXT NOT NULL,
	    user_type VARCHAR(50) NOT NULL DEFAULT 'client' CHECK (status IN ('client', 'moderator'));
	)
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = conn.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS flats(
	    id SERIAL PRIMARY KEY,
	    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
	    house_id INTEGER NOT NULL UNIQUE REFERENCES houses(id),
	    number_fl INTEGER NOT NULL,
	    price INTEGER NOT NULL CHECK (price >= 0),
	    rooms INTEGER NOT NULL CHECK (rooms >= 1)),
	    status VARCHAR(50) NOT NULL CHECK (status IN ('created', 'approved', 'declined', 'on moderation'));
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: conn}, nil
}

//CREATE OR REPLACE FUNCTION update_last_flat_added()
//RETURNS TRIGGER AS $$
//BEGIN
//UPDATE houses
//SET last_flat_added = CURRENT_TIMESTAMP
//WHERE id = NEW.house_id;
//RETURN NEW;
//END;
//$$ LANGUAGE plpgsql;

//CREATE TRIGGER update_last_apartment_added_trigger
//AFTER INSERT ON flats
//FOR EACH ROW
//EXECUTE FUNCTION update_last_apartment_added();

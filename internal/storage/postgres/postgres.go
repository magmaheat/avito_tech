package postgres

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Storage struct {
	db *pgxpool.Pool
}

type House struct {
	id        int64
	address   string
	year      int64
	developer interface{}
	createdFl time.Time
}

type Flat struct {
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.postgres.New"
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	// TODO add sql injections
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS houses (
			id INTEGER PRIMARY KEY,
			address TEXT NOT NULL,
			year INTEGER NOT NULL CHECK (year >= 0),
			developer TEXT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			update_at TIMESTAMP NULL
		);
	`)

	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		user_type VARCHAR(50) NOT NULL DEFAULT 'client' CHECK (user_type IN ('client', 'moderator'))
		);
	`)

	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS flats (
			id SERIAL PRIMARY KEY,
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			house_id INTEGER NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
			number_flat INTEGER NOT NULL,
			price INTEGER NOT NULL CHECK (price >= 0),
			rooms INTEGER NOT NULL CHECK (rooms >= 1),
			status VARCHAR(50) NOT NULL CHECK (status IN ('created', 'approved', 'declined', 'on moderation'))
		);
	`)

	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: pool}, nil
}

func (s *Storage) Create(id, year int64, address, developer string) error {
	const fn = "storage.postgres.CreateHouse"

	var developerValue interface{}
	if developer == "" {
		developerValue = nil
	} else {
		developerValue = developer
	}

	house := House{
		address:   address,
		year:      year,
		developer: developerValue,
		createdFl: time.Now(),
	}

	query, args, err := squirrel.
		Insert("houses").
		Columns("address", "year", "developer", "created_at").
		Values(house.address, house.year, house.developer, time.Now()).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	ctx := context.Background()
	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

//func (s *Storage) Get(id int64) (int64, error) {
//
//}

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

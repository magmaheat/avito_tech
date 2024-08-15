package postgres

import (
	"avito_tech/internal/entity"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Storage struct {
	db *pgxpool.Pool
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
			house_id INTEGER NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
			number INTEGER NOT NULL,
			price INTEGER NOT NULL CHECK (price >= 0),
			rooms INTEGER NOT NULL CHECK (rooms >= 1),
			status VARCHAR(50) NOT NULL CHECK (status IN ('created', 'approved', 'declined', 'on moderation'))
		);
	`)

	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = pool.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_flats_house_id_status
		ON flats (house_id, status);
	`)

	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: pool}, nil
}

func (s *Storage) CreateHouse(house entity.House) error {
	const fn = "storage.postgres.CreateHouse"

	var developerValue interface{}
	if house.Developer == "" {
		developerValue = nil
	} else {
		developerValue = house.Developer
	}

	query, args, err := squirrel.
		Insert("houses").
		Columns("id", "address", "year", "developer", "created_at").
		Values(house.ID, house.Address, house.Year, developerValue, time.Now()).
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

func (s *Storage) GetFlats(id int64, role string) ([]entity.Flat, error) {
	const fn = "storage.postgres.Get"
	ctx := context.Background()

	var rows pgx.Rows
	var err error

	if role != "moderator" {
		rows, err = s.db.Query(ctx, `
			SELECT house_id, number, price, rooms
			FROM flats
			WHERE house_id = $1 and status = 'approved'
		`, id)

	} else {
		rows, err = s.db.Query(ctx, `
			SELECT *
			FROM flats
			WHERE house_id = $1
		`, id)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var flats []entity.Flat

	for rows.Next() {
		var flat entity.Flat
		err := rows.Scan(
			&flat.HouseID,
			&flat.Number,
			&flat.Price,
			&flat.Rooms,
			&flat.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		flats = append(flats, flat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return flats, nil
}

func (s *Storage) CreateFlat(flat entity.Flat) error {
	const fn = "storage.postgres.CreateFlag"

	query, args, err := squirrel.
		Insert("flats").
		Columns("house_id", "number", "price", "rooms", "status").
		Values(flat.HouseID, flat.Number, flat.Price, flat.Rooms, "created").
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

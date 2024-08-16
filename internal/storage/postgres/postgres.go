package postgres

import (
	"avito_tech/internal/entity"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
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
			user_id UUID NOT NULL REFERENCES users(id),
			house_id INTEGER NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
			number INTEGER NOT NULL CHECK (number >= 1),
			price INTEGER NOT NULL CHECK (price >= 0),
			rooms INTEGER NOT NULL CHECK (rooms >= 1),
			status VARCHAR(50) NOT NULL CHECK (status IN ('created', 'approved', 'declined', 'on moderation')),
			last_moderator_id UUID NULL
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

	_, err = pool.Exec(ctx, `
		CREATE OR REPLACE FUNCTION func_update_at()
		RETURNS TRIGGER AS $$
		BEGIN
			UPDATE houses
			SET update_at = CURRENT_TIMESTAMP
			WHERE id = NEW.house_id;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = pool.Exec(ctx, `
		CREATE TRIGGER func_update_at_trigger
		AFTER INSERT ON flats
		FOR EACH ROW
		EXECUTE FUNCTION func_update_at();
	`)

	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: pool}, nil
}

func (s *Storage) CreateUser(user entity.User, hashPassword string) (uuid.UUID, error) {
	const fn = "storage.postgres.CreateUser"

	query, args, err := squirrel.
		Insert("users").
		Columns("email", "password", "user_type").
		Values(user.Email, hashPassword, user.UserType).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", fn, err)
	}

	var id uuid.UUID
	ctx := context.Background()

	err = s.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func (s *Storage) CreateH(house entity.House) (int64, error) {
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
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return -1, fmt.Errorf("%s: %w", fn, err)
	}

	var id int64
	ctx := context.Background()

	err = s.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func (s *Storage) GetFlats(id int64, role string) ([]entity.Flat, error) {
	const fn = "storage.postgres.Get"
	ctx := context.Background()

	var rows pgx.Rows
	var err error

	if role != "moderator" {
		rows, err = s.db.Query(ctx, `
			SELECT *
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

func (s *Storage) CreateF(flat entity.Flat) (int64, error) {
	const fn = "storage.postgres.CreateFlag"

	query, args, err := squirrel.
		Insert("flats").
		Columns("user_id", "house_id", "number", "price", "rooms", "status").
		Values(flat.UserID, flat.HouseID, flat.Number, flat.Price, flat.Rooms, "created").
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return -1, fmt.Errorf("%s: %w", fn, err)
	}

	var id int64
	ctx := context.Background()

	err = s.db.QueryRow(ctx, query, args...).Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func (s *Storage) Update(flat entity.Flat, idMod uuid.UUID) error {
	const fn = "storage.postgres.Update"

	if !checkFlat(flat) {
		return fmt.Errorf("invalid arguments: %s", fn)
	}

	queryBuilder := squirrel.Update("flats").
		Where(squirrel.And{
			squirrel.Eq{"id": flat.ID},
			squirrel.Or{
				squirrel.NotEq{"status": "on moderation"},
				squirrel.Eq{"last_moderator_id": idMod},
			},
		}).
		Set("house_id", flat.HouseID).
		Set("number", flat.Number).
		Set("price", flat.Price).
		Set("rooms", flat.Rooms).
		Set("status", flat.Status).
		Set("last_moderator_id", idMod).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query %s: %w", fn, err)
	}

	res, err := s.db.Exec(context.Background(), query, args...)
	if err != nil {
		return fmt.Errorf("failed to update flat %s: %w", fn, err)
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were update %s", fn)
	}

	return nil
}

func checkFlat(flat entity.Flat) bool {
	if flat.ID == 0 ||
		flat.HouseID == 0 ||
		flat.Number == 0 ||
		flat.Price == 0 ||
		flat.Rooms == 0 ||
		flat.Status == "" {
		return false
	}

	return true
}

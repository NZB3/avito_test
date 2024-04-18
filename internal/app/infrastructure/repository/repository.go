package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
	"project/internal/app/models"
	"project/internal/logger"
)

type repository struct {
	db  *sql.DB
	log logger.Logger
}

func New(log logger.Logger) (*repository, error) {
	const op = "repository.New"
	cfg, err := loadConfig()
	if err != nil {
		log.Errorf("%s Failed to load database config: %s", op, err)
		return nil, err
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Errorf("%s Failed to connect to database: %s", op, err)
		return nil, err
	}

	return &repository{
		log: log,
		db:  db,
	}, nil
}

func (r *repository) GetBanner(ctx context.Context, tagID int, featureID int) (models.Banner, error) {
	const op = "repository.GetBanner"

	row := r.db.QueryRowContext(ctx, `SELECT id, tag_ids, feature_id, content, is_active, created_at, updated_at
FROM banners WHERE $1=ANY(tag_ids) AND feature_id = $2`, tagID, featureID)
	if err := row.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Banner{}, models.BannerNotFound
		}
		r.log.Errorf("%s Failed to execute query: %s", op, err)
		return models.Banner{}, err
	}

	var bannerDB dbBanner
	err := row.Scan(
		&bannerDB.ID,
		pq.Array(&bannerDB.TagIDs),
		&bannerDB.FeatureID,
		&bannerDB.Content,
		&bannerDB.IsActive,
		&bannerDB.CreatedAt,
		&bannerDB.UpdatedAt,
	)

	if err != nil {
		r.log.Errorf("%s Failed to scan row: %s", op, err)
		return models.Banner{}, err
	}

	banner := mapOnBanner(bannerDB)
	return banner, nil
}

func (r *repository) GetBanners(ctx context.Context, limit, offset int) ([]models.Banner, error) {
	const op = "repository.GetBanners"
	rows, err := r.db.QueryContext(ctx, `SELECT id, tag_ids, feature_id, content, is_active, created_at, updated_at
FROM banners ORDER BY created_at LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		r.log.Errorf("%s Failed to execute query: %s", op, err)
		return nil, err
	}
	defer rows.Close()

	banners := make([]models.Banner, 0)
	for rows.Next() {
		var bannerDB dbBanner
		err := rows.Scan(
			&bannerDB.ID,
			pq.Array(&bannerDB.TagIDs),
			&bannerDB.FeatureID,
			&bannerDB.Content,
			&bannerDB.IsActive,
			&bannerDB.CreatedAt,
			&bannerDB.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("%s Failed to scan row: %s", op, err)
			return nil, err
		}

		banner := mapOnBanner(bannerDB)
		banners = append(banners, banner)
	}

	return banners, nil
}

func (r *repository) GetBannersByTagID(ctx context.Context, tagID int, limit int, offset int) ([]models.Banner, error) {
	const op = "repository.GetBannersByTagID"

	rows, err := r.db.QueryContext(ctx, `SELECT id, tag_ids, feature_id, content, is_active, created_at, updated_at
FROM banners WHERE $1=ANY(tag_ids) ORDER BY created_at LIMIT $2 OFFSET $3`, tagID, limit, offset)
	if err != nil {
		r.log.Errorf("%s Failed to execute query: %s", op, err)
		return nil, err
	}
	defer rows.Close()

	banners := make([]models.Banner, 0)
	for rows.Next() {
		var bannerDB dbBanner
		err := rows.Scan(
			&bannerDB.ID,
			pq.Array(&bannerDB.TagIDs),
			&bannerDB.FeatureID,
			&bannerDB.Content,
			&bannerDB.IsActive,
			&bannerDB.CreatedAt,
			&bannerDB.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("%s Failed to scan row: %s", op, err)
		}
		banner := mapOnBanner(bannerDB)
		banners = append(banners, banner)
	}

	return banners, nil
}

func (r *repository) GetBannersByFeatureID(ctx context.Context, featureID int, limit int, offset int) ([]models.Banner, error) {
	const op = "repository.GetBannersByFeatureID"

	rows, err := r.db.QueryContext(ctx, `SELECT id, tag_ids, feature_id, content, is_active, created_at, updated_at
FROM banners WHERE feature_id=$1 ORDER BY created_at LIMIT $2 OFFSET $3`, featureID, limit, offset)
	if err != nil {
		r.log.Errorf("%s Failed to execute query: %s", op, err)
		return nil, err
	}
	defer rows.Close()

	banners := make([]models.Banner, 0)
	for rows.Next() {
		var bannerDB dbBanner
		err := rows.Scan(
			&bannerDB.ID,
			pq.Array(&bannerDB.TagIDs),
			&bannerDB.FeatureID,
			&bannerDB.Content,
			&bannerDB.IsActive,
			&bannerDB.CreatedAt,
			&bannerDB.UpdatedAt,
		)
		if err != nil {
			r.log.Errorf("%s Failed to scan row: %s", op, err)
		}
		banner := mapOnBanner(bannerDB)
		banners = append(banners, banner)
	}

	return banners, nil
}

func (r *repository) UpdateBanner(ctx context.Context, banner models.Banner) (bool, error) {
	const op = "repository.UpdateBanner"

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.log.Errorf("%s Failed to begin transaction: %s", op, err)
		return false, err
	}
	defer tx.Rollback()

	bannerDB := mapOnDBBanner(banner)
	res, err := tx.ExecContext(ctx, `UPDATE banners SET tag_ids=$1, feature_id=$2, content=$3, is_active=$4 WHERE id=$5`,
		pq.Array(bannerDB.TagIDs), bannerDB.FeatureID, bannerDB.Content, bannerDB.IsActive, bannerDB.ID)
	if err != nil {
		r.log.Errorf("%s Failed to execute query: %s", op, err)
		return false, err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		r.log.Errorf("%s Failed to get rows affected: %s", op, err)
		return false, err
	}

	return affectedRows != 0, nil
}

func (r *repository) CreateBanner(ctx context.Context, banner models.Banner) (int, error) {
	const op = "repository.CreateBanner"

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.log.Errorf("%s Failed to begin transaction: %s", op, err)
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO banners (tag_ids, feature_id, content, is_active) values ($1, $2, $3, $4) RETURNING id`)
	if err != nil {
		r.log.Errorf("%s Failed to prepare query: %s", op, err)
		return 0, err
	}
	defer stmt.Close()

	var id int
	bannerDB := mapOnDBBanner(banner)
	err = stmt.QueryRowContext(ctx, pq.Array(bannerDB.TagIDs), bannerDB.FeatureID, bannerDB.Content, bannerDB.IsActive).Scan(&id)
	if err != nil {
		r.log.Errorf("%s Failed to get last insert ID: %s", op, err)
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		r.log.Errorf("%s Failed to commit transaction: %s", op, err)
		return 0, err
	}

	return id, nil
}

func (r *repository) DeleteBanner(ctx context.Context, bannerID int) (bool, error) {
	const op = "repository.DeleteBanner"

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.log.Errorf("%s Failed to begin transaction: %s", op, err)
		return false, err
	}
	defer tx.Rollback()

	res, err := r.db.ExecContext(ctx, `DELETE FROM banners WHERE id=$1`, bannerID)
	if err != nil {
		r.log.Errorf("%s Failed to execute query: %s", op, err)
		return false, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		r.log.Errorf("%s Failed to get rows affected: %s", op, err)
		return false, err
	}

	return affected != 0, nil
}

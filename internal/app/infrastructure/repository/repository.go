package repository

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)
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

func (r *repository) GetBanner(ctx context.Context, tagID int, featureID int) (*models.Banner, error) {
	const op = "repository.GetBanner"

	row := r.db.QueryRowContext(ctx, `SELECT (id, tag_ids, feature_id, content, is_active, created_at, updated_at) 
FROM banners WHERE $1=ANY(tag_ids) AND feature_id = $2`, tagID, featureID)
	if err := row.Err(); err != nil {
		r.log.Errorf("%s Failed to execute query: %s", op, err)
		return nil, err
	}

	var banner models.Banner
	err := row.Scan(&banner)
	if err != nil {
		r.log.Errorf("%s Failed to scan into models.Banner: %s", op, err)
		return nil, err
	}

	return &banner, nil
}

func (r *repository) GetBannersByTag(ctx context.Context, tagID int, limit int, offset int) ([]*models.Banner, error) {
	const op = "repository.GetBannersByTag"

	rows, err := r.db.QueryContext(ctx, `SELECT (id, tag_ids, feature_id, content, is_active, created_at, updated_at)
FROM banners WHERE $1=ANY(tag_ids) ORDER BY created_at LIMIT $2 OFFSET $3`, tagID, limit, offset)
	if err != nil {
		r.log.Errorf("%s Failed to execute query: %s", op, err)
		return nil, err
	}
	defer rows.Close()

	var banners []*models.Banner
	for rows.Next() {
		var banner models.Banner
		err := rows.Scan(&banner)
		if err != nil {
			r.log.Errorf("%s Failed to : %s", op, err)
			return nil, err
		}

		banners = append(banners, &banner)
	}

	return banners, nil
}

func (r *repository) GetBannersByFeatureID(ctx context.Context, featureID int, limit int, offset int) ([]*models.Banner, error) {
	const op = "repository.GetBannersByFeatureID"

	rows, err := r.db.QueryContext(ctx, `SELECT (id, tag_ids, feature_id, content, is_active, created_at, updated_at)
FROM banners WHERE feature_id=$1`, featureID)
	if err != nil {
		r.log.Errorf("%s Failed to execute query: %s", op, err)
		return nil, err
	}
	defer rows.Close()

	var banners []*models.Banner
	for rows.Next() {
		var banner models.Banner
		err := rows.Scan(&banner)
		if err != nil {
			r.log.Errorf("%s Failed to : %s", op, err)
		}
	}

	return banners, nil
}

func (r *repository) UpdateBanner(ctx context.Context, banner models.Banner) (bool, error) {
	const op = "repository.UpdateBanner"

	res, err := r.db.ExecContext(ctx, `UPDATE banners SET tag_ids=$1, feature_id=$2, content=$3, is_active=$4 WHERE id=$5`,
		banner.TagIDs, banner.FeatureID, banner.Content, banner.IsActive, banner.ID)
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

func (r *repository) CreateBanner(ctx context.Context, banner models.Banner) (int, error) {
	const op = "repository.CreateBanner"
	res, err := r.db.ExecContext(ctx, `INSERT INTO banners values (tag_ids, feature_id, content, is_active) RETURNING id`, banner.TagIDs, banner.FeatureID, banner.Content, banner.IsActive)
	if err != nil {
		r.log.Errorf("%s Failed to execute query: %s", op, err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		r.log.Errorf("%s Failed to get last insert ID: %s", op, err)
		return 0, err
	}

	return int(id), nil
}

func (r *repository) DeleteBanner(ctx context.Context, bannerID int) (bool, error) {
	const op = "repository.DeleteBanner"

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

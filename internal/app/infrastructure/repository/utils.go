package repository

import (
	"encoding/json"
	"project/internal/app/models"
	"time"
)

type dbBanner struct {
	ID        int       `db:"id"`
	TagIDs    []int32   `db:"tag_ids"`
	FeatureID int       `db:"feature_id"`
	Content   []byte    `db:"content"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func mapOnDBBanner(banner models.Banner) dbBanner {
	content, err := json.Marshal(&banner.Content)
	if err != nil {
		panic(err)
	}

	tagIDs := make([]int32, len(banner.TagIDs))
	for i, tagID := range banner.TagIDs {
		tagIDs[i] = int32(tagID)
	}

	return dbBanner{
		ID:        banner.ID,
		TagIDs:    tagIDs,
		FeatureID: banner.FeatureID,
		Content:   content,
		IsActive:  banner.IsActive,
		CreatedAt: banner.CreatedAt,
	}
}

func mapOnBanner(bannerDB dbBanner) models.Banner {
	var content map[string]interface{}
	err := json.Unmarshal(bannerDB.Content, &content)
	if err != nil {
		panic(err)
	}

	tagIDs := make([]int, len(bannerDB.TagIDs))
	for i, tagID := range bannerDB.TagIDs {
		tagIDs[i] = int(tagID)
	}

	return models.Banner{
		ID:        bannerDB.ID,
		TagIDs:    tagIDs,
		FeatureID: bannerDB.FeatureID,
		Content:   content,
		IsActive:  bannerDB.IsActive,
		CreatedAt: bannerDB.CreatedAt,
		UpdatedAt: bannerDB.UpdatedAt,
	}
}

package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

type Banner struct {
	ID        int       `json:"banner_id"`
	TagIDs    []int     `json:"tag_ids"`
	FeatureID int       `json:"feature_id"`
	Content   []byte    `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *Banner) TagIDsFeatureIDHash() string {
	hashString := fmt.Sprintf("%v%v", b.TagIDs, b.FeatureID)
	hash := md5.Sum([]byte(hashString))
	hashHex := hex.EncodeToString(hash[:])
	return hashHex
}

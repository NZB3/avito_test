package models

type User struct {
	ID     uint64 `json:"id"`
	TagIDs []int  `json:"tag_ids"`
}

package admin

import "time"

type Stats struct {
	Users   int32 `json:"users"`
	Heroes  int32 `json:"heroes"`
	Images  int32 `json:"images"`
	Storage int64 `json:"storage"`
}

type UserSummary struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
}

type ContainerSummary struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

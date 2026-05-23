package admin

type Stats struct {
	Users   int32 `json:"users"`
	Heroes  int32 `json:"heroes"`
	Images  int32 `json:"images"`
	Storage int64 `json:"storage"`
}

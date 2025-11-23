package domain

type RequestSign struct {
	Fullname *string `json:"fullname"`
	Email    string  `json:"email" binding:"required"`
	Pw       string  `json:"pw" binding:"required"`
}

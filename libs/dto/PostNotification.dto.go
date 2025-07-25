package dto

type PostNotificationDTO struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

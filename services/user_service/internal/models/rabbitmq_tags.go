package models

type ChangeTagsTutorToBroker struct {
	TutorTelegramID int64    `json:"tutor_telegram_id"`
	Tags            []string `json:"tags"`
}

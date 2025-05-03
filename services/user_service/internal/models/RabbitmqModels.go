package models

type AddResponsesToTutor struct {
	TutorTelegramID int64 `json:"tutor_telegram_id"`
	ResponseCount   int   `json:"response_count"`
}

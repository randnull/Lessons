package models

import "time"

type Review struct {
	ID        string    `json:"id"`
	TutorID   string    `json:"tutor_id"`
	StudentID string    `json:"student_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type ReviewRequest struct {
	TutorID string `json:"tutor_id"`
	Comment string `json:"comment"`
	Rating  int    `json:"rating"`
}

package models

import "time"

type CreateUser struct {
	Name       string `json:"name"`
	TelegramId int64  `json:"telegram_id"`
}

type TutorsPagination struct {
	Tutors []*TutorForList
	Pages  int
}

type User struct {
	Id         string `json:"id"`
	TelegramID int64  `json:"telegram_id"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	IsBanned   bool   `json:"is_banned"`
}

type Tutor struct {
	Id         string
	TelegramID int64
	Name       string
	Bio        string
	Tags       []string
}

type TutorForList struct {
	Id     string
	Name   string
	Rating int32
	Tags   []string
}

type TutorDetails struct {
	Tutor         User
	Bio           string
	ResponseCount int32
	Reviews       []Review
	IsActive      bool
	Tags          []string
	Rating        int32
	CreatedAt     time.Time
}

type UpdateBioTutor struct {
	Bio string `json:"bio"`
}

type UpdateTagsTutor struct {
	Tags []string `json:"tags"`
}
type ChangeTagsTutorToBroker struct {
	TutorTelegramID int64    `json:"tutor_telegram_id"`
	Tags            []string `json:"tags"`
}

type ChangeActive struct {
	IsActive bool `json:"is_active"`
}

type UpdateNameTutor struct {
	Name string `json:"name"`
}

type BanUser struct {
	TelegramID int64 `json:"telegram_id"`
	IsBan      bool  `json:"is_ban"`
}

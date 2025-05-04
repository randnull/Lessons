package models

import "time"

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
	Bio string `json:"bio" validate:"required,min=10,max=1000"`
}

type UpdateTagsTutor struct {
	Tags []string `json:"tags" validate:"required,min=1,max=10,dive,required"`
}

type UpdateNameTutor struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type ChangeTagsTutorToBroker struct {
	TutorTelegramID int64    `json:"tutor_telegram_id"`
	Tags            []string `json:"tags"`
}

type BanUser struct {
	TelegramID int64 `json:"telegram_id" validate:"required,gt=0"`
	IsBan      bool  `json:"is_ban" validate:"required"`
}

package models

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
}

type Tutor struct {
	Id         string
	TelegramID int64
	Name       string
	Bio        string
	Tags       []string
}

type TutorForList struct {
	Id   string
	Name string
	Tags []string
}

type TutorDetails struct {
	Tutor         User
	Bio           string
	ResponseCount int32
	Reviews       []Review
	Tags          []string
}

type UpdateBioTutor struct {
	Bio string `json:"bio"`
}

type UpdateTagsTutor struct {
	Tags []string `json:"tags"`
}

type ChangeActive struct {
	IsActive bool `json:"is_active"`
}

package models

type CreateUser struct {
	Name       string `json:"name"`
	TelegramId int64  `json:"telegram_id"`
}

type TutorsPagination struct {
	User  []*User
	Pages int
}

type User struct {
	Id         string `json:"id"`
	TelegramID int64  `json:"telegram_id"`
	Name       string `json:"name"`
	//CreatedAt  time.Time `json:"created_at"`
}

type Tutor struct {
	Id string `json:"id"`
	//TelegramID int64  `json:"telegram_id"`
	Bio  string `json:"bio"`
	Name string `json:"name"`
	//CreatedAt  time.Time `json:"created_at"`
}

type TutorModel struct {
	User User
	Bio  string `json:"bio"`
}

type TutorDetails struct {
	Reviews []Review
	Tags    []string
	Tutor   TutorModel
	Bio     string
}

type UpdateBioTutor struct {
	Bio string `json:"bio"`
}

type UpdateTagsTutor struct {
	Tags []string `json:"tags"`
}

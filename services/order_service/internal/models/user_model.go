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

type UpdateBioTutor struct {
	Bio string `json:"bio"`
}

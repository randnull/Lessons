package models

type User struct {
	UserId  string `json:"user_id"`
	Univer  string `json:"univer"`
	AboutMe string `json:"info"`
}

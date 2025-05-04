package models

type OrderToBroker struct {
	ID        string   `json:"order_id"`
	StudentID int64    `json:"student_id"`
	Title     string   `json:"order_name"`
	Tags      []string `json:"tags"`
	Status    string   `json:"status"`
}

type SuggestOrder struct {
	ID              string `json:"order_id"`
	TutorTelegramID int64  `json:"tutor_telegram_id"`
	Title           string `json:"order_name"`
	Description     string `json:"description"`
	MinPrice        int    `json:"min_price"`
	MaxPrice        int    `json:"max_price"`
}

type ReviewToBroker struct {
	ReviewID        string `json:"review_id"`
	ResponseID      string `json:"response_id"`
	OrderID         string `json:"order_id"`
	OrderName       string `json:"order_name"`
	TutorTelegramID int64  `json:"tutor_telegram_id"`
}

type ResponseToBrokerModel struct {
	ResponseID string `json:"response_id"`
	TutorID    int64  `json:"tutor_id"`
	StudentID  int64  `json:"student_id"`
	OrderID    string `json:"order_id"`
	Title      string `json:"order_name"`
}

type ChangeTagsTutorToBroker struct {
	TutorTelegramID int64    `json:"tutor_telegram_id"`
	Tags            []string `json:"tags"`
}

type SelectedResponseToBroker struct {
	OrderID    string `json:"order_id"`
	OrderName  string `json:"order_name"`
	StudentID  int64  `json:"student_telegram_id"`
	TutorID    int64  `json:"tutor_telegram_id"`
	ResponseID string `json:"response_id"`
}

package models

type RequestGetUserById struct {
	Id int `json:"id"`
}

type RequestSignUp struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RequestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RequestCreateCategory struct {
	Name string `json:"name"`
}

type RequestGetCategories struct {
	Name string `json:"q"`
	RequestPagination
}

type RequestGetCategoryById struct {
	Id int `json:"id" uri:"id"`
}

type RequestUpdateCategory struct {
	Name string `json:"name"`
}

type RequestCreateTransaction struct {
	Amount     float64 `json:"amount"`
	Type       string  `json:"type"`
	CategoryId int     `json:"category_id"`
}

type RequestGetTransactions struct {
	UserId     int    `json:"user_id"`
	CategoryId string `json:"category_id"`
	Type       string `json:"type"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	RequestPagination
}

type RequestGetTransactionById struct {
	Id int `json:"id" uri:"id"`
}

type RequestUpdateTransaction struct {
	Amount     float64 `json:"amount"`
	Type       string  `json:"type"`
	CategoryId string  `json:"category_id"`
}

type QueryPagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Page   int `json:"page"`
}

type RequestPagination struct {
	Limit string `json:"limit"`
	Page  string `json:"page"`
}

type RequestGetBalance struct {
	UserId    int    `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

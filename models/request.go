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

type QueryPagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type RequestPagination struct {
	Limit string `json:"limit"`
	Page  string `json:"page"`
}


package models

type Response struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseLogin struct {
	User
	Token string `json:"token"`
}

type ResponseCategoryList struct {
	Data []Category `json:"data"`
	Count    int64      `json:"count"`
	Page     int        `json:"page"`
	Limit    int        `json:"limit"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type LoginResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

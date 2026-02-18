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
	Data  []Category `json:"data"`
	Count int64      `json:"count"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
}

type ResponseTransactionList struct {
	Data  []TransactionResponse `json:"data"`
	Count int64                 `json:"count"`
	Page  int                   `json:"page"`
	Limit int                   `json:"limit"`
}

type TransactionResponse struct {
	Id        int                    `json:"id"`
	User      UserSimpleResponse     `json:"user"`
	Amount    float64                `json:"amount"`
	Type      string                 `json:"type"`
	Category  CategorySimpleResponse `json:"category"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}

type UserSimpleResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CategorySimpleResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type ResponseBalance struct {
	UserId       int     `json:"user_id"`
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
}

type ResponseUserList struct {
	Data  []UserResponse `json:"data"`
	Count int64          `json:"count"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

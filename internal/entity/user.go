package entity

type User struct {
	Id       int    `json:"-" db:"id"`
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

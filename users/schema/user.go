package schema

type User struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

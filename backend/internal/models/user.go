package models

type User struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	HashedPassword string `json:"hashed_password"`
}

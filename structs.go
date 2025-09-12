package main

import "time"

type Todo struct {
	Id          int        `json:"id" gorm:"column:id"`
	Title       string     `json:"title" binding:"required"`
	Completed   *bool      `json:"completed" binding:"required"`
	Category    string     `json:"category" binding:"required"`
	Priority    string     `json:"priority" binding:"required"`
	CompletedAt *time.Time `json:"completedAt" gorm:"column:completedat"`
	DueDate     *time.Time `json:"dueDate" gorm:"column:duedate"`
}

type TodoDTO struct {
	Completed *bool `json:"completed" binding:"required"`
}

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Username string `json:"username" binding:"required" gorm:"unique"`
	Password string `json:"password" binding:"required" gorm:"not null"` //bcypt hash + required
	Role     string `json:"role" gorm:"default:user"`
}

type UserDTO struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func ToUserDTO(user User) UserDTO {
	return UserDTO{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
}

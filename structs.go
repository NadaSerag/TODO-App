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

type User struct{
   	ID    uint   `gorm:"primaryKey"`
    Username  string `gorm:"unique"`
    Password string `gorm:"size:255;not null"` //bcypt hash + required
		Role string `gorm:"default:user"`
}

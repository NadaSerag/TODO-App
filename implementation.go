package main

//file to implement the APIs

import (
	"fmt"
	"strconv"
	"time"

	"github.com/NadaSerag/TODO-App/middleware"
	"github.com/gin-gonic/gin"
)

// connectionStr is a formatted string that tells Go’s PostgreSQL driver how to connect to our database.
//var connectionStr = "postgres://postgres:Mydatabase123@localhost:5432/todo_app?sslmode=disable"

//var db, _ = sql.Open("postgres", connectionStr)

// - 200: Successful GET, POST, PUT, or DELETE.
// - 400: Invalid JSON or empty title.
// - 404: Todo not found for GET, PUT, or DELETE.

// When a request comes in, Gin creates a *gin.Context object
// and passes it into the function as the parameter
// the parameter is by convention named "c".

func GetTodos(c *gin.Context) {
	allTodos := []Todo{}

	//Find() returns a *gorm.DB object (which contains things like error status, rows affected, etc.).
	result := DB.Find(&allTodos) // SELECT * FROM users;

	//code 500: for Internal server/database error.
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error})
	}

	//just used these lines to make sure everything is going fine:
	//fmt.Println(result.RowsAffected)  correct, printed 7, and they are 7 todos
	//fmt.Println(allTodos)             correct, printed all todos in terminal

	//200 code for Successful GET
	c.JSON(200, allTodos)

}

func GetTodoById(c *gin.Context) {
	//c.Param(...) returns the value as a string.
	//thats wy we need to convert it to an int by strconv.Atoi
	id, _ := strconv.Atoi(c.Param("id"))
	var todoByIdFound Todo

	result := DB.Find(&todoByIdFound, id)

	//why not " if result.Error != nil " ?
	//beacuse updating/deleting/getting a row that doesn’t exist is not considered an error!!
	//fmt.Println(result.Error)  -> Prints <nil>

	if result.RowsAffected == 0 {
		// code 404 used bec: Todo not found for GET, PUT, or DELETE.
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}

	//row successfully found
	c.JSON(200, todoByIdFound)
}

func GetTodosByCategory(c *gin.Context) {
	//detecting the wanted category in the URL
	category := c.Param("category")

	//the returning response initialized as an empty array of todos
	todosByCategory := []Todo{}

	if !CategoryValid(category) {
		c.JSON(400, gin.H{"error": "Invalid Category in URL: Category must be Work, Personal, or Study"})
		return
	}

	result := DB.Where("category = ?", category).Find(&todosByCategory)

	if result.RowsAffected == 0 {
		// code 404 used bec: Todo not found for GET, PUT, or DELETE.
		c.JSON(404, gin.H{"error": fmt.Sprintf("No todos with category = %s", category)})
		return
	}

	c.JSON(200, todosByCategory)
}

func GetTodosByStatus(c *gin.Context) {
	//detecting the wanted status from the URL
	stat := c.Param("status")

	//the returning response initialized as an empty array of todos
	todosByStatus := []Todo{}

	result := DB.Where("completed = ?", stat).Find(&todosByStatus)

	if result.RowsAffected == 0 {
		// code 404 for Resource Not Found
		c.JSON(404, gin.H{"error": fmt.Sprintf("No todos with status = %s", stat)})
		return
	}

	c.JSON(200, todosByStatus)

}

func GetTodosBySearch(c *gin.Context) {
	//detecting the wanted status from the URL
	search := c.Query("q")

	//the returning response initialized as an empty array of todos
	todosBySearch := []Todo{}

	result := DB.Where("title LIKE ?", search).Find(&todosBySearch)

	if result.RowsAffected == 0 {
		// code 404 used bec: Todo not found for GET, PUT, or DELETE.
		c.JSON(404, gin.H{"error": fmt.Sprintf("No todos have titles that include  '%s'", search)})
		return
	}
	c.JSON(200, todosBySearch)
}

func CreateTodo(c *gin.Context) {
	var newTodo Todo

	//VERY IMPORTANT: reading the request body (JSON)

	//N.B.:
	//c.Param(" ... ") = extracts directly from the URL, not the JSON body.
	//BindJSON or ShouldBindJson  = reads the JSON payload, not the URL.

	//returns an error, error is = nil if JSON parses successfully
	err := c.ShouldBindJSON(&newTodo)

	gottenClaims, exists_ok := middleware.GetUserClaims(c)

	everythingOk := middleware.ClaimsCheck(c, gottenClaims, exists_ok)

	
	//if false if returned by ClaimsCheck, then it c.aborted in the function ClaimsCheck already, so we just exit (return)
	if !everythingOk {
		return
	}

	//checking if the title is empty
	if newTodo.Title == "" {
		c.JSON(400, gin.H{"error": "Title cannot be empty"})
		return
	}

	if !PriorityValid(newTodo.Priority) {
		c.JSON(400, gin.H{"error": "Priority must be Low, Medium, or High"})
		return
	}

	if !CategoryValid(newTodo.Category) {
		c.JSON(400, gin.H{"error": "Category must be Work, Personal, or Study"})
		return
	}

	if !DueDateValid(newTodo.DueDate) {
		c.JSON(400, gin.H{"error": "Due date cannot be in the past"})
		return
	}

	//checking for invalid JSON
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	//SETTING COMPLETION STAT
	if newTodo.Completed != nil && *(newTodo.Completed) {
		//placing the time.Now in a variable (type time.Time)
		currentTime := time.Now()
		//taking its address &currentTime to assign it to CompletedAt (pointer to time.Time: type *time.Time)
		newTodo.CompletedAt = &currentTime
	}

	//SETTING USER_ID
	newTodo.UserID = gottenClaims.UserID

	DB.Create(&newTodo)
	c.JSON(200, newTodo)
}

func UpdateTodo(c *gin.Context) {

	var updatedTodo Todo
	err := c.ShouldBindJSON(&updatedTodo)

	if updatedTodo.Title == "" {
		c.JSON(400, gin.H{"error": "Title cannot be empty"})
		return
	}

	if !PriorityValid(updatedTodo.Priority) {
		c.JSON(400, gin.H{"error": "Priority must be Low, Medium, or High"})
		return
	}

	if !CategoryValid(updatedTodo.Category) {
		c.JSON(400, gin.H{"error": "Category must be Work, Personal, or Study"})
		return
	}

	if !DueDateValid(updatedTodo.DueDate) {
		c.JSON(400, gin.H{"error": "Due date cannot be in the past"})
		return
	}

	//checking for invalid JSON
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	idToSearchFor, _ := strconv.Atoi(c.Param("id"))
	var giveID Todo
	giveID.Id = idToSearchFor

	result := DB.Model(&giveID).Updates(updatedTodo)

	// if result.Error != nil {
	// 	c.JSON(404, gin.H{"error": "Todo not found, probably invalid ID"})
	// 	return
	// }
	//fmt.Println(result.Error)

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Todo not found, probably invalid ID"})
		return
	}
	//why not " if result.Error != nil " ?
	//beacuse updating/deleting/getting a row that doesn’t exist is not considered an error
	//fmt.Println(result.Error)  -> Prints <nil>

	if updatedTodo.Completed != nil && *(updatedTodo.Completed) {
		//placing the time.Now in a variable (type time.Time)
		currentTime := time.Now()
		//taking its address &currentTime to assign it to CompletedAt (pointer to time.Time: type *time.Time)
		updatedTodo.CompletedAt = &currentTime

		DB.Model(&Todo{}).Where("id = ?", idToSearchFor).Update("completedat", currentTime)
	}

	if updatedTodo.Completed != nil && !*(updatedTodo.Completed) {
		DB.Model(&Todo{}).Where("id = ?", idToSearchFor).Update("completedat", nil)
	}
	updatedTodo.Id = idToSearchFor

	c.JSON(200, updatedTodo)
}

func UpdateTodosByCategory(c *gin.Context) {
	var updatedTodos []Todo
	var updatedStat TodoDTO

	categoryToSearchFor := c.Param("category")

	if !CategoryValid(categoryToSearchFor) {
		c.JSON(400, gin.H{"error": "Invalid Category in URL: Category must be Work, Personal, or Study"})
		return
	}

	err := c.ShouldBindJSON(&updatedStat)

	//checking for invalid JSON
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	DB.Model(&Todo{}).Where("category = ?", categoryToSearchFor).Update("completed", updatedStat.Completed)

	if updatedStat.Completed != nil && *(updatedStat.Completed) {
		currentTime := time.Now()
		DB.Model(&Todo{}).Where("category = ?", categoryToSearchFor).Update("completedat", currentTime)
	}

	if updatedStat.Completed != nil && !*(updatedStat.Completed) {
		DB.Model(&Todo{}).Where("category = ?", categoryToSearchFor).Update("completedat", nil)
	}

	DB.Where("category = ?", categoryToSearchFor).Find(&updatedTodos)

	//200 code for Successful PUT
	c.JSON(200, updatedTodos)
}

func DeleteById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	result := DB.Delete(&Todo{}, id)

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Todo not found, probably invalid ID"})
		return
	}

	//why not " if result.Error != nil " ?
	//beacuse updating/deleting/getting a row that doesn’t exist is not considered an error
	//fmt.Println(result.Error)  -> Prints <nil>

	//200 code for Successful DELETE
	c.JSON(200, gin.H{"message": "Todo deleted"})
}

func DeleteAll(c *gin.Context) {

	//deleting a;; rows, soft deletion
	DB.Where("1 = 1").Delete(&Todo{})

	//200 code for Successful DELETE
	c.JSON(200, gin.H{"message": "All todos deleted"})
}

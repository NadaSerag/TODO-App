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

// GetTodos returns list of all todos
// @Summary      List all todos
// @Description  Returns a list of all todos
// @Tags         Todos
// @Success      200  {array} Todo
// @Produce json
// @Router       /todos [get]
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

type ErrorJSON struct {
	ErrorStr string `json:"error"`
}

type SuccessJSON struct {
	SuccessStr string `json:"message"`
}

// GetTodoById returns a single todo by ID
// @Summary Get a todo by ID
// @Description Returns a single todo by its ID
// @Tags Todos
// @Success 200 {object} Todo
// @Failure 404 {object} ErrorJSON "Not Found"
// @Param        id   path      int  true  "Todo ID"
// @Router /todos/{id} [get]
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

// GetTodoByCategory returns a single todo by category
// @Summary Get a todo by category
// @Description Returns a single todo by its category
// @Tags Todos
// @Success 200 {array} Todo
// @Failure 404 {object} ErrorJSON "Category Invalid"
// @Param        category   path      string  true  "Category to search for: "
// @Router /category/{category} [get]
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

// GetTodoByStatus returns a single todo by status
// @Summary Get a todo by status
// @Description Returns a single todo by its status
// @Tags Todos
// @Success 200 {array} Todo
// @Failure 404 {object} ErrorJSON "No todos with the status"
// @Param        status   path      string  true  "Status to search for: "
// @Router /status/{status} [get]
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

// GetTodoBySearch returns todos that have name-match with the search string entered
// @Summary Search Todos by title
// @Description Returns todos that have name-match with the search string entered
// @Tags Todos
// @Success 200 {array} Todo
// @Failure 404 {object} ErrorJSON "No todos match search"
// @Param        q   query      string  true  "search string: "
// @Router /todos/search/ [get]
func GetTodosBySearch(c *gin.Context) {
	//detecting the wanted status from the URL
	search := c.Query("q")

	//the returning response initialized as an empty array of todos
	todosBySearch := []Todo{}

	result := DB.Where("title LIKE ?", "%"+search+"%").Find(&todosBySearch)

	if result.RowsAffected == 0 {
		// code 404 used bec: Todo not found for GET, PUT, or DELETE.
		c.JSON(404, gin.H{"error": fmt.Sprintf("No todos have titles that include  '%s'", search)})
		return
	}
	c.JSON(200, todosBySearch)
}

// CreateTodo adds a new todo to our table
// @Summary Add todo
// @Description Adding/Creating new todo to out table "todos"
// @Tags Todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        todo  body      Todo  true  "Todo to create"
// @Success      201   {object}  Todo
// @Failure      400   {object}  ErrorJSON "Invalid json, or missing fields, or invalid priority/category or past due date"
// @Failure      401   {object}  ErrorJSON "Unauthorized"
// @Router       /todos [post]
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

// @Summary Edit a todo
// @Description Editing a todo by entering its id, and the request includes the updated data
// @Tags Todos
// @Accept       json
// @Produce      json
// Security     BearerAuth
// @Param        todo  body      Todo  true  "Updates"
// @Success      201   {object}  Todo
// @Failure      400   {object}  ErrorJSON "Invalid json, or missing fields, or invalid priority/category or past due date"
// @Failure      401   {object}  ErrorJSON "Invalid Token"
// @Failure      403   {object}  ErrorJSON "Unauthorized (Forbidden for that person to do that action)"
// @Router       /todos/{id} [put]
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

// @Summary Edit todos
// @Description Editing multiple todos of the same category to update their state of completion
// @Tags Todos
// @Accept       json
// @Produce      json
// Security     BearerAuth
// @Param        completion_status  body      TodoDTO  true  "Updates"
// @Success      201   {object}  Todo
// @Failure      400   {object}  ErrorJSON "Invalid json"
// Failure      401   {object}  ErrorJSON "Unauthorized"
// @Router       /todos/category/{category}} [put]
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

// @Summary Delete a todo
// @Description Deleting a todo by its id
// @Tags Todos
// Security     BearerAuth
// @Param        id   path      int  true  "Todo ID to delete"
// @Success      201   {object}  SuccessJSON "Todo deleted successfully"
// @Failure      400   {object}  ErrorJSON "No todo exists with the ID entered"
// Failure      401   {object}  ErrorJSON "Unauthorized"
// @Router       /todos/{id} [delete]
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

// @Summary Delete ALL todos
// @Description Deleting ALL todos in the entire database - this action is only permissible for an Admin.
// @Tags Todos
// @Security     BearerAuth
// @Success      201   {object}  SuccessJSON "All todos deleted"
// @Failure      400   {object}  ErrorJSON "Invalid Token"
// @Failure      401   {object}  ErrorJSON "Unauthorized"
// @Router       /todos [delete]
func DeleteAll(c *gin.Context) {

	//deleting a;; rows, soft deletion
	DB.Where("1 = 1").Delete(&Todo{})

	//200 code for Successful DELETE
	c.JSON(200, gin.H{"message": "All todos deleted"})
}

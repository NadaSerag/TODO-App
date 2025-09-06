package main

//file to implement the APIs

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Id          int        `json:"id"`
	Title       string     `json:"title" binding:"required"`
	Completed   *bool      `json:"completed" binding:"required"`
	Category    string     `json:"category" binding:"required"`
	Priority    string     `json:"priority" binding:"required"`
	CompletedAt *time.Time `json:"completedAt"`
	DueDate     *time.Time `json:"dueDate"`
}

// connection with DSN (Data Source Name) established
// connectionStr is a formatted string that tells Go’s PostgreSQL driver how to connect to our database.
var connectionStr = "postgres://postgres:Mydatabase123@localhost:5432/todo_app?sslmode=disable"

var db, _ = sql.Open("postgres", connectionStr)

// var Todos []Todo
var id = 5

// - 200: Successful GET, POST, PUT, or DELETE.
// - 400: Invalid JSON or empty title.
// - 404: Todo not found for GET, PUT, or DELETE.

// When a request comes in, Gin creates a *gin.Context object
// and passes it into the function as the parameter
// the parameter is by convention named "c".

func GetTodos(c *gin.Context) {
	allTodos := []Todo{}
	var todoGotten Todo

	//retrieving todos from our database to return them
	query := `
    SELECT *
    FROM todos
`
	rows, _ := db.Query(query)
	for rows.Next() {
		rows.Scan(
			&todoGotten.Id,
			&todoGotten.Title,
			&todoGotten.Completed,
			&todoGotten.Category,
			&todoGotten.Priority,
			&todoGotten.CompletedAt,
			&todoGotten.DueDate,
		)
		allTodos = append(allTodos, todoGotten)
	}

	//200 code for Successful GET
	c.JSON(200, allTodos)

	rows.Close()

}

func GetTodoById(c *gin.Context) {
	//c.Param(...) returns the value as a string.
	//thats wy we need to convert it to an int by strconv.Atoi
	id, _ := strconv.Atoi(c.Param("id"))
	var todoByIdFound Todo

	//method#1: searching in out in-memory struct
	// for i := 0; i < len(Todos); i++ {
	// 	if Todos[i].Id == id {
	// 		c.JSON(200, Todos[i])
	// 		return
	// 	}
	// }

	//method#2 ( better ): searching in our database
	query := `
    SELECT *
    FROM todos
    WHERE id = $1
`
	row := db.QueryRow(query, id)

	err := row.Scan(
		&todoByIdFound.Id,
		&todoByIdFound.Title,
		&todoByIdFound.Completed,
		&todoByIdFound.Category,
		&todoByIdFound.Priority,
		&todoByIdFound.CompletedAt,
		&todoByIdFound.DueDate,
	)

	if err != nil {
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
	var todoToAdd Todo

	//by using our struct
	// for i := 0; i < len(Todos); i++ {
	// 	if Todos[i].Category == category {
	// 		todosByCategory = append(todosByCategory, Todos[i])
	// 	}
	// }

	//by using our database & SQL instructions
	query := `
    SELECT *
    FROM todos
    WHERE category = $1
`
	rows, _ := db.Query(query, category)
	for rows.Next() {
		rows.Scan(
			&todoToAdd.Id,
			&todoToAdd.Title,
			&todoToAdd.Completed,
			&todoToAdd.Category,
			&todoToAdd.Priority,
			&todoToAdd.CompletedAt,
			&todoToAdd.DueDate,
		)
		todosByCategory = append(todosByCategory, todoToAdd)
	}
	c.JSON(200, todosByCategory)
	rows.Close()

	//LOOP EXPLANATION (for ...)

	// 	rows.Scan importance:
	// Postgres gives Go a row of raw values (like a slice of interfaces — []interface{}),
	// but Go has no idea how I want to map them into my struct.
	// rows.Scan takes those raw column values and copies them into the variables in order.

	//rows.Next usage:
	// the rows object represents a stream of results coming from the database.
	// rows.Next() advances to the next row in the result set.
	// It returns true if there is another row, false if there are no more.
	//Thats why it's usually called inside a for loop.
}

func GetTodosByStatus(c *gin.Context) {
	//detecting the wanted status from the URL
	stat := c.Param("status")

	//the returning response initialized as an empty array of todos
	todosByStatus := []Todo{}
	var todoToAdd Todo

	query := `
    SELECT *
    FROM todos
    WHERE completed = $1
`
	rows, _ := db.Query(query, stat)
	for rows.Next() {
		rows.Scan(
			&todoToAdd.Id,
			&todoToAdd.Title,
			&todoToAdd.Completed,
			&todoToAdd.Category,
			&todoToAdd.Priority,
			&todoToAdd.CompletedAt,
			&todoToAdd.DueDate,
		)
		todosByStatus = append(todosByStatus, todoToAdd)
	}
	c.JSON(200, todosByStatus)
	rows.Close()
}

func GetTodosBySearch(c *gin.Context) {
	//detecting the wanted status from the URL
	search := c.Param("search")

	//the returning response initialized as an empty array of todos
	todosBySearch := []Todo{}
	var todoToAdd Todo

	query := `
    SELECT *
    FROM todos
    WHERE title LIKE $1
`
	rows, _ := db.Query(query, search+"%")

	for rows.Next() {
		rows.Scan(
			&todoToAdd.Id,
			&todoToAdd.Title,
			&todoToAdd.Completed,
			&todoToAdd.Category,
			&todoToAdd.Priority,
			&todoToAdd.CompletedAt,
			&todoToAdd.DueDate,
		)
		todosBySearch = append(todosBySearch, todoToAdd)
	}
	c.JSON(200, todosBySearch)
	rows.Close()
}

func CreateTodo(c *gin.Context) {
	var newTodo Todo

	//VERY IMPORTANT: reading the request body (JSON)

	//N.B.:
	//c.Param(" ... ") = extracts directly from the URL, not the JSON body.
	//BindJSON or ShouldBindJson  = reads the JSON payload, not the URL.

	//returns an error, error is = nil if JSON parses successfully
	err := c.ShouldBindJSON(&newTodo)

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

	query := `
        INSERT INTO todos (title, completed, category, priority, completedAt, dueDate)
        VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING id
    `

	db.QueryRow(query,
		newTodo.Title,
		newTodo.Completed,
		newTodo.Category,
		newTodo.Priority,
		newTodo.CompletedAt,
		newTodo.DueDate,
	).Scan(&newTodo.Id)

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

	query := `
        UPDATE todos
				SET title = $1, completed = $2, category = $3, priority = $4, dueDate = $5
				WHERE id = $6
				`

	idToSearchFor, _ := strconv.Atoi(c.Param("id"))

	res, _ := db.Exec(query,
		updatedTodo.Title,
		updatedTodo.Completed,
		updatedTodo.Category,
		updatedTodo.Priority,
		updatedTodo.DueDate,
		idToSearchFor)

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Todo not found, probably invalid ID"})
		return
	}

	if updatedTodo.Completed != nil && *(updatedTodo.Completed) {
		//placing the time.Now in a variable (type time.Time)
		currentTime := time.Now()
		//taking its address &currentTime to assign it to CompletedAt (pointer to time.Time: type *time.Time)
		updatedTodo.CompletedAt = &currentTime
	}

	updatedTodo.Id = idToSearchFor

	c.JSON(200, updatedTodo)

}

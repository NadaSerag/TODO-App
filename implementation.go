package main

//file to implement the APIs

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Id          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Completed   *bool  `json:"completed" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Priority    string `json:"priority" binding:"required"`
	CompletedAt string `json:"completedAt"`
	DueDate     string `json:"dueDate" binding:"required"`
}

// connection with DSN (Data Source Name) established
// connectionStr is a formatted string that tells Go’s PostgreSQL driver how to connect to our database.
var connectionStr = "postgres://postgres:Mydatabase123@localhost:5432/todo_app?sslmode=disable"

var db, _ = sql.Open("postgres", connectionStr)
var Todos []Todo
var id = 1

func InitializeTodoArray() []Todo {
	//The in-memory store (slice) is initialized as empty.
	Todos = []Todo{}
	return Todos
}

// - 200: Successful GET, POST, PUT, or DELETE.
// - 400: Invalid JSON or empty title.
// - 404: Todo not found for GET, PUT, or DELETE.

// When a request comes in, Gin creates a *gin.Context object
// and passes it into the function as the parameter
// the parameter is by convention named "c".

func GetTodos(c *gin.Context) {
	//200 code for Successful GET
	c.JSON(200, Todos)
}

func GetTodoById(c *gin.Context) {
	//c.Param(...) returns the value as a string.
	//thats wy we need to convert it to an int by strconv.Atoi
	id, _ := strconv.Atoi(c.Param("id"))

	for i := 0; i < len(Todos); i++ {
		if Todos[i].Id == id {
			c.JSON(200, Todos[i])
			return
		}
	}

	// code 404 used bec: Todo not found for GET, PUT, or DELETE.
	c.JSON(404, gin.H{"error": "Todo not found"})
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
	errorVar := c.ShouldBindJSON(&newTodo)

	//checking for invalid JSON
	if errorVar != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	//checking if the title is empty
	if newTodo.Title == "" {
		c.JSON(400, gin.H{"error": "Title cannot be empty"})
		return
	}

	newTodo.Id = id
	Todos = append(Todos, newTodo)
	//incrementing manually?
	id++

	query := `
        INSERT INTO todos (id, title, completed, category, priority, completedAt, dueDate)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	db.Exec(query,
		newTodo.Id,
		newTodo.Title,
		newTodo.Completed,
		newTodo.Category,
		newTodo.Priority,
		newTodo.CompletedAt,
		newTodo.DueDate,
	)
	c.JSON(200, newTodo)
}

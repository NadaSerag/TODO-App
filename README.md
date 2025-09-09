# TODO-App
A RESTful API for a Todo App using Go and the `gin-gonic/gin` framework for routing.

The app manages todo items stored in a PostgreSQL database using GORM (ORM for Go) and supports CRUD operations (Create, Read, Update, Delete) with features like categorization, prioritization, completion tracking, due dates, status filtering, title search, and bulk updates.
### Database set up:
1- Download PostgreSQL (make sure to memorize or write the password you enter while installing because it's needed in step 4)

2- Create your database (name it as you like)

3- Don't worry about creating a table for now, the app already creates table 'todos' on setup if it doesn't exist

4- Edit the DSN connection in main.go near line 29 (the connectionStr variable)
- connectionStr := "postgres://postgres:YOUR-PASSWORD@localhost:5432/YOUR-DATABASE-NAME?sslmode=disable"
- switch YOUR-PASSWORD with the password you entered while installing
- & switch YOUR-DATABASE-NAME with the name you named your database earlier in step 2.

5- That's It! Set up completed.

### Recommended: How to Test Each Endpoint
Here are the steps on how to start testing each endpoint. Trust me, it's so much fun when you get the hang of it!

1- Download Postman

2- Select the method (GET, POST, PUT, DELETE)

3- Type in the URL

4- Enter the request body if needed (in JSON)

### Running the App with Docker

1- Install Docker Desktop

2- Build and start the containers by running "docker-compose up --build" in the terminal

3- The API will be available at:
   http://localhost:8080

4- To stop the containers run "docker-compose down"


#### There is a Test Evidence Document provided among the files of this repository, containing screenshots of successful testing of every endpoint. Both the .pdf or .docx versions are available for downloading.




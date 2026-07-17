# Go Swagger API (Employee Information Service)

A RESTful API built with **Go (Golang)** that provides endpoints to fetch employee information from a **MySQL** database. This project also serves as a foundational exploration of Go language capabilities and integrates **Swagger** for interactive API documentation.

---

## Local Development Setup

Follow these steps to set up the environment and run the application locally.

### 1. Prerequisites
Ensure you have the following installed on your machine:
* **Go** (1.16+ recommended)
* **MySQL Server**

### 2. Environment Configuration
The application requires specific environment variables to connect to the database and start the server. 

Create a `.env` file in the root directory and populate it with your local configurations using this structure:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
API_HOST=localhost
API_PORT=8080
ALLOWED_ORIGINS=*
```

### 3. Install Dependencies
Download and tidy up the required Go modules and packages:
```bash
go mod tidy
```

### 4. Run the Application
Start the Go API server:
```bash
go run main.go
```

---

## API Documentation (Swagger)
Once the server is running, you can explore, test, and interact with the endpoints via the built-in Swagger UI. 

Open your browser and navigate to:
```text
http://localhost:8080/swagger/index.html
```
*(Note: Adjust the port above if you changed `API_PORT` in your `.env` file).*

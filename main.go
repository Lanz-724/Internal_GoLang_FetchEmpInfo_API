package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"

	docs "my-api/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           My Go API
// @version         1.0
// @description     My first Go API
// @host            localhost:9090
// @BasePath        /
var db *sql.DB

func main() {
	var err error

	// Use = instead of := since err is already declared
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read from env
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	apiPort := os.Getenv("API_PORT")

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	origins := strings.Split(allowedOrigins, ",")

	host := os.Getenv("API_HOST") + ":" + os.Getenv("API_PORT")
	docs.SwaggerInfo.Host = host

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening DB:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}
	fmt.Println("Connected to MySQL!")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/users", getUsers)
	r.GET("/empinfo", getEmpInfo)
	r.Run(":" + apiPort)
}

// Employee struct maps to your table columns
type Employee struct {
	ID    int    `json:"id"`
	Lname string `json:"lname"`
	Fname string `json:"fname"`
}

// @Summary      Get employees
// @Description  Returns a list of employees
// @Tags         employees
// @Produce      json
// @Param        limit   query  int  false  "Number of records to return" default(20)
// @Param        offset  query  int  false  "Number of records to skip"   default(0)
// @Success      200  {array}   Employee
// @Router       /users [get]
func getUsers(c *gin.Context) {
	// Get params from URL, with default values if not provided
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")

	query := fmt.Sprintf("SELECT empi_idno, empi_lname, empi_fname FROM tks_employee_info LIMIT %s, %s", offset, limit)

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var emp Employee
		if err := rows.Scan(&emp.ID, &emp.Lname, &emp.Fname); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		employees = append(employees, emp)
	}

	c.JSON(http.StatusOK, employees)
}

type EmpInfo struct {
	Name  string `json:"Name"`
	Dept  string `json:"Dept"`
	IdNo  string `json:"IdNo"`
	BioID string `json:"BioID"`
}

type EmpInfoResponse struct {
	Total     int       `json:"total"`
	Employees []EmpInfo `json:"employees"`
}

// @Summary      Get employee info
// @Description  Returns employees with department and bio info
// @Tags         employees
// @Produce      json
// @Param        Dept   query  string  false  "Department name to filter by" default()
// @Param        FilterName   query  string  false  "Filter by Emp Name to filter by" default()
// @Success      200  {array}   EmpInfo
// @Router       /empinfo [get]
func getEmpInfo(c *gin.Context) {

	dept := c.DefaultQuery("Dept", "")
	filterName := c.DefaultQuery("FilterName", "")

	// Add % wildcards for LIKE
	deptFilter := "%" + dept + "%"
	nameFilter := "%" + filterName + "%"

	query := `SELECT
    CONCAT(UCASE(info.empi_lname), ' ', UCASE(info.empi_fname)) AS NAME,
		dept.di_name AS Dept,
		info.empi_idno AS IdNo,
		bio.ebi_uid AS BioID
	FROM tks_employee_info info
		INNER JOIN sys_dept_info dept ON info.di_id = dept.di_id
		INNER JOIN tks_emp_bioidentity bio ON info.empi_id = bio.empi_id
	WHERE 
		dept.di_name LIKE ?
		AND (info.empi_lname LIKE ? OR info.empi_fname LIKE ?)
		
	ORDER BY NAME asc`

	rows, err := db.Query(query, deptFilter, nameFilter, nameFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var employees []EmpInfo
	for rows.Next() {
		var emp EmpInfo
		if err := rows.Scan(&emp.Name, &emp.Dept, &emp.IdNo, &emp.BioID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		employees = append(employees, emp)
	}

	response := EmpInfoResponse{
		Total:     len(employees),
		Employees: employees,
	}

	c.JSON(http.StatusOK, response)

}

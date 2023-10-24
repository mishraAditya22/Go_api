package Controllers

import (
	"api/Config"
	"api/Models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	Db := Config.ConnectDatabase()
	defer Db.Close()
	if Db == nil {
		return
	}
	User := Models.User{}
	err := c.ShouldBindJSON(&User)
	if err != nil {
		fmt.Println("Error while mapping user inputs !")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Bad Input Variables !"})
		return
	}
	Query := `insert into userData (username,password) values($1,$2)`
	_, errQ := Db.Exec(Query, User.Username, User.Password)
	if errQ != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Couldn't create the new user."})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "User is successfully created."})
	}
}

func GetUserByID(c *gin.Context) {
	Db := Config.ConnectDatabase()
	defer Db.Close()
	if Db == nil {
		return
	}

	userID := c.Param("id")

	// Check if the userID is empty
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID is missing from the request"})
		return
	}

	// Convert the userID to an integer
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid User ID"})
		return
	}

	// Query the database to retrieve the user by ID
	var user Models.User
	query := "SELECT * FROM userData WHERE id = $1"
	row := Db.QueryRow(query, id)
	err = row.Scan(&user.Id, &user.Username, &user.Password)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)

}

func GetUsers(c *gin.Context) {
	Db := Config.ConnectDatabase()
	if Db == nil {
		return
	}
	Query := "select * from userData"
	rows, err := Db.Query(Query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error while Quering the Data-Base , try again later !"})
		return
	}
	var Users []Models.User
	for rows.Next() {
		var user Models.User
		row_Err := rows.Scan(&user.Id, &user.Username, &user.Password)
		if row_Err != nil {
			c.AbortWithStatusJSON(http.StatusInsufficientStorage, gin.H{"message": "Query Failed !"})
			return
		}
		Users = append(Users, user)
	}
	c.IndentedJSON(http.StatusOK, Users)
}

func UpdateUser(c *gin.Context) {
	Db := Config.ConnectDatabase()
	if Db == nil {
		return
	}
	defer Db.Close()
	userID := c.Param("id")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User ID is missing from the request"})
		return
	}
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid User ID"})
		return
	}

	var UserInput Models.User
	Err := c.ShouldBindJSON(&UserInput)
	if Err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid User-Inputs !"})
		return
	}
	Query := "update userData set username = $1 where id=$2"
	row, Err := Db.Exec(Query, UserInput.Username, id)
	if Err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error while executing Database Query !"})
		return
	}
	c.IndentedJSON(http.StatusOK, row)
}

func DeleteUser(c *gin.Context) {
	Db := Config.ConnectDatabase()
	if Db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to connect to the database"})
		return
	}

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID is missing from the request"})
		Db.Close()
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid User ID"})
		Db.Close()
		return
	}

	query := "DELETE FROM userData WHERE id = $1"
	_, err = Db.Exec(query, id)
	if err != nil {
		// Log the error for debugging
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while executing Database Query"})
		Db.Close()
		return
	}

	Db.Close()
	c.IndentedJSON(http.StatusOK, gin.H{"message": "User record deleted with User-Id: " + userID})
}

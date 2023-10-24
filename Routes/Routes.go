package Routes

import (
	"api/Controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {

	R1 := router.Group("/user")

	//To create a new user
	R1.POST("/create", Controllers.CreateUser)
	//To get users by their id
	R1.GET("/:id", Controllers.GetUserByID)
	//To get all the users in the system
	R1.GET("/", Controllers.GetUsers)
	//To Update user's data
	R1.PUT("/:id", Controllers.UpdateUser)
	//To delete User's Record
	R1.DELETE("/:id", Controllers.DeleteUser)
}

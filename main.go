package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    // "log"
    "fmt"
    // "context"
)

type User struct {
    ID     string  `json:"id"`
    Firstname  string  `json:"firstname"`
    Lastname string  `json:"lastname"`
}

var users = []User{
    {ID: "1", Firstname: "Colby", Lastname: "Chenard"},
    {ID: "2", Firstname: "Jenny", Lastname: "Wallace"},
    {ID: "3", Firstname: "Dutch", Lastname: "eroni"},
}

func main() {
    // manager := NewMongoDBManager()
	// client, err := manager.ConnectToMongoDB()
	// if err != nil {
	// 	log.Fatalf("Error connecting to MongoDB: %v", err)
	// }
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		log.Fatalf("Failed to disconnect MongoDB client: %v", err)
	// 	}
	// 	fmt.Println("Disconnected from MongoDB successfully")
	// }()

    router := gin.Default()

    // Kasa routes
    router.GET("/device/action/:id", getUsers)
    router.POST("/device/:id/:action", action)

    // User routes
    router.GET("/users/:id", getUserByID)
    router.PUT("/users/:id", updateUser)
    // router.PUT("/users/create/:id", createUser)

    router.POST("/users/create", func(c *gin.Context) {
        var user User
        c.BindJSON(&user)
        fmt.Printf("User to create: %v\n", user)
        c.JSON(http.StatusOK, gin.H{"user": user.Firstname})

    })
        router.DELETE("/users/:id", deleteUser)
    
    // Hue routes

    // Start the server
    router.Run(":8080")
}





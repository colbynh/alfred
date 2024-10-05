package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "os/exec"
    "log"
)

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// In-memory user storage
var users = []User{}

func main() {
    router := gin.Default()

    // Routes
    router.GET("/device/action/:id", getUsers)
    router.POST("/device/:id/:action", action)
    router.GET("/users/:id", getUserByID)
    router.PUT("/users/:id", updateUser)
    router.DELETE("/users/:id", deleteUser)

    // Start the server
    router.Run(":8080")
}

func getUsers(c *gin.Context) {
    c.JSON(http.StatusOK, users)
}

func action(c *gin.Context) {
	id := c.Param("id")
    action := c.Param("action")

    // Construct the command
    cmd := exec.Command("kasa", "--host", id, action)
    output, err := cmd.Output()

    // Handle error in command execution
    if err != nil {
        log.Printf("Error executing command: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute command", "details": err.Error()})
        return
    }

    // Return output as JSON response
    c.JSON(http.StatusOK, gin.H{"message": "Command executed successfully", "output": string(output)})
}

func getUserByID(c *gin.Context) {
    id := c.Param("id")
    for _, user := range users {
        if user.ID == id {
            c.JSON(http.StatusOK, user)
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

func updateUser(c *gin.Context) {
    id := c.Param("id")
    var updatedUser User
    if err := c.BindJSON(&updatedUser); err != nil {
        return
    }

    for i, user := range users {
        if user.ID == id {
            users[i] = updatedUser
            c.JSON(http.StatusOK, updatedUser)
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

func deleteUser(c *gin.Context) {
    id := c.Param("id")
    for i, user := range users {
        if user.ID == id {
            users = append(users[:i], users[i+1:]...)
            c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

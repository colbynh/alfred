package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
	// "os/exec"
    // "log"
    // "fmt"
)


// type User struct {
// 	First         string        `bson:"first,omitempty"`
// 	Last          string        `bson:"last,omitempty"`
// }


// func createUser(c *gin.Context) {
//     coll := client.Database("homeauto").Collection("user")
//     newUser := Restaurant{Name: "Colby", Last: "Bolby"}

//     result, err := coll.InsertOne(context.TODO(), newUser)
//     if err != nil {
//         fmt.Println(err)
//         c.JSON(http.StatusNotFound, gin.H{"message": fmt.Printf("Failed to create user: %v\n", err)})
//     }
//     msg := fmt.Printf("New user %user")
//     c.JSON(http.StatusOK, gin.H{"message": "New user created"})
// }

func getUsers(c *gin.Context) {
    c.JSON(http.StatusOK, users)
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


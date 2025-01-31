package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
	"os/exec"
    "log"
    "fmt"
)

func action(c *gin.Context) {
	id := c.Param("id")
    action := c.Param("action")

    // Construct the command
    cmd := exec.Command("kasa", "--host", id, action)
    output, err := cmd.Output()
    fmt.Println(output)
    // Handle error in command execution
    if err != nil {
        log.Printf("Error executing command: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute command", "details": err.Error()})
        return
    }

    // Return output as JSON response
    c.JSON(http.StatusOK, gin.H{"message": "Command executed successfully", "output": string(output)})
}
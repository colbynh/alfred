package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestAction(t *testing.T) {
    // Set Gin to Test mode
    gin.SetMode(gin.TestMode)

    // Create a new router
    router := gin.Default()
    router.POST("/device/:id/:action", action)

    // Create a test request with a valid ID and action
    req, err := http.NewRequest(http.MethodPost, "/device/192.168.101.170", "on")
    if err != nil {
        t.Fatalf("Could not create request: %v", err)
    }

    // Record the response
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)

    // Check the response code
    assert.Equal(t, http.StatusOK, recorder.Code)

    // Check the response body contains a success message
    expectedBody := "{\"message\":\"Command executed successfully\"}"
    assert.Contains(t, recorder.Body.String(), expectedBody)
} 
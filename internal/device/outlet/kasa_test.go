package outlet

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// Change from var to a pointer
var ScanOpenPortsTest = &scanOpenPortsFunc{
	f: func() ([]string, []error) {
		return nil, nil
	},
}

type scanOpenPortsFunc struct {
	f func() ([]string, []error)
}

type mock struct {
	discoverDevicesIps func() (map[string]interface{}, error)
}

func (m *mock) DiscoverDevicesIps() (map[string]interface{}, error) {
	return m.discoverDevicesIps()
}

// Add this at package level
var scanOpenPortsFn = func() ([]string, []error) {
	return nil, nil
}

func TestDiscoverDevicesIps(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	k := &kasaOutlet{
		id:     "test-id",
		logger: logger,
		c:      &gin.Context{},
	}

	jsonData, err := k.discoverDevicesIps()
	assert.NoError(t, err)
	assert.NotNil(t, jsonData)

	// Check if "ips" key exists
	ipsInterface, exists := jsonData["ips"]
	assert.True(t, exists, "ips key should exist in jsonData")

	// Convert interface{} to []interface{} first
	ipsSlice, ok := ipsInterface.([]interface{})
	assert.True(t, ok, "ips should be a slice")

	// Convert each element to string and check
	var ips []string
	for _, ip := range ipsSlice {
		if strIP, ok := ip.(string); ok {
			ips = append(ips, strIP)
		}
	}

	assert.Contains(t, ips, "192.168.101.43")
	assert.Contains(t, ips, "192.168.101.170")
}

func TestState(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	k := &kasaOutlet{
		id:     "test-id",
		logger: logger,
	}

	// Mock the exec.Command function
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "Device state: True")
	}

	jsonData, err := k.state()
	assert.NoError(t, err)
	assert.NotNil(t, jsonData)
	assert.Equal(t, "True", jsonData["state"])
}

func TestSysInfo(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	k := &kasaOutlet{
		id:     "test-id",
		logger: logger,
	}

	// Mock the exec.Command function
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", `{"model": "HS103(US)", "sw_ver": "1.0.13"}`)
	}

	jsonData, err := k.sysInfo()
	assert.NoError(t, err)
	assert.NotNil(t, jsonData)
	assert.Equal(t, "HS103(US)", jsonData["model"])
	assert.Equal(t, "1.0.13", jsonData["sw_ver"])
}

func TestAction(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	k := &kasaOutlet{
		id:     "test-id",
		logger: logger,
	}

	// Mock the exec.Command function
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "OK")
	}

	router := gin.Default()
	router.POST("/api/v1/device/outlet/:brand/:id/:action", func(c *gin.Context) {
		action := c.Param("action")
		c.Param("id")
		c.Param("brand")
		err := k.action(action, c)
		assert.NoError(t, err)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/device/outlet/kasa/192.168.101.170/off", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
}

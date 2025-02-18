package outlet

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"
	"time"

	"net"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type mock struct {
	discoverDevicesIps func() (map[string]interface{}, error)
}

func (m *mock) DiscoverDevicesIps() (map[string]interface{}, error) {
	return m.discoverDevicesIps()
}

// Add at package level
type dialTimeoutFunc struct {
	f func(network, addr string, timeout time.Duration) (net.Conn, error)
}

var dialTimeoutWrapper = &dialTimeoutFunc{
	f: func(network, addr string, timeout time.Duration) (net.Conn, error) {
		return nil, nil
	},
}

// Add at package level, near the top of the file
type scanOpenPortsFunc struct {
	f func() ([]string, []error)
}

var scanOpenPortsWrapper = &scanOpenPortsFunc{
	f: func() ([]string, []error) {
		return nil, nil
	},
}

func TestDiscoverDevicesIps(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	k := &kasaOutlet{
		id:     "test-id",
		logger: logger,
		c:      &gin.Context{},
	}

	// Mock ScanOpenPorts to return consistent test data
	originalScan := scanOpenPortsWrapper.f
	defer func() { scanOpenPortsWrapper.f = originalScan }()

	scanOpenPortsWrapper.f = func() ([]string, []error) {
		return []string{"192.168.101.170"}, nil
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

func TestScanOpenPorts(t *testing.T) {
	// Store original function
	originalDialTimeout := dialTimeoutWrapper.f
	defer func() { dialTimeoutWrapper.f = originalDialTimeout }()

	// Test error case first
	dialTimeoutWrapper.f = func(network, addr string, timeout time.Duration) (net.Conn, error) {
		return nil, errors.New("connection failed")
	}

	// Mock the scan function
	originalScan := scanOpenPortsWrapper.f
	defer func() { scanOpenPortsWrapper.f = originalScan }()

	scanOpenPortsWrapper.f = func() ([]string, []error) {
		return nil, []error{errors.New("no open ports found")}
	}

	ips, errs := scanOpenPortsWrapper.f()
	assert.NotEmpty(t, errs)
	assert.Empty(t, ips)
	assert.Equal(t, "no open ports found", errs[0].Error())

	// Then test successful case
	dialTimeoutWrapper.f = func(network, addr string, timeout time.Duration) (net.Conn, error) {
		return &mockConn{}, nil
	}

	scanOpenPortsWrapper.f = func() ([]string, []error) {
		return []string{"192.168.101.170"}, nil
	}

	ips, errs = scanOpenPortsWrapper.f()
	assert.Empty(t, errs)
	assert.NotEmpty(t, ips)
	assert.Contains(t, ips, "192.168.101.170")
}

// mockConn implements net.Conn interface with minimal implementation
type mockConn struct{}

func (m *mockConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (m *mockConn) Write(b []byte) (n int, err error)  { return 0, nil }
func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

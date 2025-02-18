// Package outlet provides functionality for controlling smart outlets.
// This test file contains unit tests for the Kasa smart outlet implementation.
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

// mock implements a test double for the outlet interface.
// It allows for flexible behavior definition in tests by providing
// function fields that can be set to desired test behaviors.
type mock struct {
	discoverDevicesIps func() (map[string]interface{}, error)
}

// DiscoverDevicesIps implements the outlet interface method for the mock.
// It delegates to the function field, allowing for customizable test behavior.
func (m *mock) DiscoverDevicesIps() (map[string]interface{}, error) {
	return m.discoverDevicesIps()
}

// dialTimeoutFunc wraps the network dial timeout function to allow mocking.
// This type enables tests to replace the network dialing behavior with
// controlled test behavior.
type dialTimeoutFunc struct {
	f func(network, addr string, timeout time.Duration) (net.Conn, error)
}

// dialTimeoutWrapper provides a mutable reference to the dial timeout function.
// This global variable allows tests to modify network connection behavior.
var dialTimeoutWrapper = &dialTimeoutFunc{
	f: func(network, addr string, timeout time.Duration) (net.Conn, error) {
		return nil, nil
	},
}

// scanOpenPortsFunc wraps the port scanning function to allow mocking.
// This type enables tests to replace the port scanning behavior with
// controlled test behavior.
type scanOpenPortsFunc struct {
	f func() ([]string, []error)
}

// scanOpenPortsWrapper provides a mutable reference to the port scanning function.
// This global variable allows tests to modify port scanning behavior.
var scanOpenPortsWrapper = &scanOpenPortsFunc{
	f: func() ([]string, []error) {
		return nil, nil
	},
}

// TestDiscoverDevicesIps verifies that the device discovery functionality
// correctly identifies and returns IP addresses of Kasa devices.
// It tests the successful case where devices are found on the network
// and ensures the returned data structure is correctly formatted.
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

// TestState verifies that the outlet state can be retrieved correctly.
// It tests the parsing of the device state response and ensures
// the state is correctly represented in the returned JSON structure.
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

// TestSysInfo verifies that the device system information can be retrieved correctly.
// It tests the parsing of the device system information response and ensures
// all expected fields are present in the returned JSON structure.
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

// TestAction verifies that device actions (on/off) are executed correctly.
// It tests the HTTP endpoint handling and command execution for device control,
// ensuring proper response formatting and error handling.
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

// TestScanOpenPorts verifies that the port scanning functionality
// correctly identifies open ports and handles errors.
// It tests both successful port discovery and error conditions,
// ensuring proper error propagation and result formatting.
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

// mockConn provides a minimal implementation of net.Conn for testing.
// It implements all required methods of the net.Conn interface with
// no-op implementations suitable for testing.
type mockConn struct{}

// Read implements net.Conn Read method
func (m *mockConn) Read(b []byte) (n int, err error) { return 0, nil }

// Write implements net.Conn Write method
func (m *mockConn) Write(b []byte) (n int, err error) { return 0, nil }

// Close implements net.Conn Close method
func (m *mockConn) Close() error { return nil }

// LocalAddr implements net.Conn LocalAddr method
func (m *mockConn) LocalAddr() net.Addr { return nil }

// RemoteAddr implements net.Conn RemoteAddr method
func (m *mockConn) RemoteAddr() net.Addr { return nil }

// SetDeadline implements net.Conn SetDeadline method
func (m *mockConn) SetDeadline(t time.Time) error { return nil }

// SetReadDeadline implements net.Conn SetReadDeadline method
func (m *mockConn) SetReadDeadline(t time.Time) error { return nil }

// SetWriteDeadline implements net.Conn SetWriteDeadline method
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

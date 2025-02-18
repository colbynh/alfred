// Package outlet provides functionality for controlling smart outlets.
// This file implements support for TP-Link Kasa smart outlets.
package outlet

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// execCommand is a package variable that wraps exec.Command for testing purposes.
var execCommand = exec.Command

// kasaOutlet represents a TP-Link Kasa smart outlet device.
// It implements the outlet interface for controlling the device state
// and retrieving device information.
type kasaOutlet struct {
	id     string         // Unique identifier (typically IP address)
	c      *gin.Context   // HTTP context for request handling
	logger *logrus.Logger // Logger for operation tracking
}

// ScanResult represents the response format for device discovery.
// It contains a list of IP addresses where Kasa devices were found.
type ScanResult struct {
	IPs []string `json:"ips"` // List of discovered device IPs
}

// Network scanning constants
const (
	timeout = 1 * time.Second // Timeout for checking each port
	port    = "9999"          // Default Kasa device port
	subnet  = "192.168.101."  // Network subnet to scan
	startIP = 1               // First IP address in scan range
	endIP   = 254             // Last IP address in scan range
)

// joinHostPort combines an IP address and port into a network address string.
func joinHostPort(ip, port string) string {
	return net.JoinHostPort(ip, port)
}

// dialTimeout attempts to establish a network connection with timeout.
func dialTimeout(network, addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

// scanIP checks if a specific IP address has the Kasa device port open.
// It is used as a goroutine in the port scanning process.
func scanIP(ip string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	address := joinHostPort(ip, port)
	conn, err := dialTimeout("tcp", address, timeout)
	if err != nil {
		return
	}
	if conn != nil {
		conn.Close()
		results <- ip
	}
}

// ScanOpenPorts scans the configured subnet for devices listening on the Kasa port.
// It returns a list of IP addresses where the port is open and any errors encountered.
func ScanOpenPorts() ([]string, error) {
	var wg sync.WaitGroup
	results := make(chan string, endIP-startIP)
	var openIPs []string

	for i := startIP; i <= endIP; i++ {
		ip := fmt.Sprintf("%s%d", subnet, i)
		wg.Add(1)
		go scanIP(ip, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for ip := range results {
		openIPs = append(openIPs, ip)
	}
	if len(openIPs) == 0 {
		return nil, errors.New("no open ports found")
	}
	return openIPs, nil
}

// getID returns the device identifier.
func (k *kasaOutlet) getID() string {
	return k.id
}

// getBrand returns the device brand name (always "kasa").
func (k *kasaOutlet) getBrand() string {
	return "kasa"
}

// discoverDevicesIps scans the network for Kasa devices and returns their IP addresses.
// The result is formatted as a JSON object with an "ips" array.
func (k *kasaOutlet) discoverDevicesIps() (map[string]interface{}, error) {
	k.logger.Debug("Scanning for open ports on subnet:", subnet)
	ips, err := ScanOpenPorts()
	if err != nil {
		k.logger.Error(err)
	}
	k.logger.Debug("Open ports found:", ips)

	sr := ScanResult{IPs: ips}
	jsonBytes, err := json.Marshal(sr)
	if err != nil {
		k.logger.Error("Error marshaling ScanResult:", err)
		return nil, err
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		k.logger.Error("Error unmarshaling JSON:", err)
		return nil, err
	}

	jsonOutput, _ := json.MarshalIndent(jsonData, "", "    ")
	k.logger.Debug("Discovered devices JSON output:", string(jsonOutput))

	return jsonData, nil
}

// state retrieves the current state (on/off) of the outlet.
// It returns the state as a JSON object with a "state" field.
func (k *kasaOutlet) state() (map[string]interface{}, error) {
	k.logger.Debug("Executing kasa state command")
	cmd := execCommand("kasa", "--host", k.id, "state")

	o, err := cmd.Output()
	if err != nil {
		k.logger.Error("Error executing kasa state command:", err)
		return nil, err
	}

	var jsonData map[string]interface{}
	re := regexp.MustCompile(`Device state:\s+(False|True)`)
	match := re.FindString(string(o))

	if match != "" {
		match = strings.ReplaceAll(match, "Device state:", "")
		jsonStr := fmt.Sprintf("{\"state\": \"%s\"}", strings.Trim(match, " "))

		if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
			k.logger.Error("Error unmarshaling JSON:", err)
			return nil, err
		}
	}
	return jsonData, nil
}

// sysInfo retrieves system information from the outlet.
// It returns device details including model and software version.
func (k *kasaOutlet) sysInfo() (map[string]interface{}, error) {
	k.logger.Debug("Executing kasa sysinfo command")
	cmd := execCommand("kasa", "--host", k.id, "sysinfo")

	o, err := cmd.Output()
	if err != nil {
		k.logger.Error("Error executing kasa sysinfo command:", err)
		return nil, err
	}

	var jsonStr string
	var jsonData map[string]interface{}
	re := regexp.MustCompile(`(?s)\{.*\}`)
	match := re.FindString(string(o))

	if match != "" {
		jsonStr = strings.ReplaceAll(match, "'", "\"")
		jsonStr = strings.ReplaceAll(jsonStr, "\"{", "{")
		jsonStr = strings.ReplaceAll(jsonStr, "}\"", "}")

		if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
			k.logger.Error("Error unmarshaling JSON:", err)
			return nil, err
		}
	}

	return jsonData, nil
}

// action executes a command on the outlet and returns the result.
// Supported actions are: "on", "off", "discover", "state", and "sysinfo".
// The result is returned as a JSON response through the gin.Context.
func (k *kasaOutlet) action(action string, c *gin.Context) error {
	k.logger.Debug("Executing action:", action)
	jsonData := map[string]interface{}{}
	var err error

	switch action {
	case "on":
		k.logger.Debug("Turning on the device")
		cmd := execCommand("kasa", "--host", k.id, "on")

		if err = cmd.Run(); err != nil {
			k.logger.Error("Error executing kasa on command:", err)
			return err
		}
	case "off":
		k.logger.Debug("Turning off the device")
		cmd := execCommand("kasa", "--host", k.id, "off")
		if err = cmd.Run(); err != nil {
			k.logger.Error("Error executing kasa off command:", err)
			return err
		}
	case "discover":
		k.logger.Debug("Discovering devices...")
		jsonData, err = k.discoverDevicesIps()
		if err != nil {
			k.logger.Error("Error discovering devices:", err)
			return err
		}
	case "state":
		k.logger.Debug("Getting device state")
		jsonData, err = k.state()
		if err != nil {
			k.logger.Error("Error getting state:", err)
			return err
		}
	case "sysinfo":
		k.logger.Debug("Getting device sysinfo")
		jsonData, err = k.sysInfo()
		if err != nil {
			k.logger.Error("Error getting sysinfo:", err)
			return err
		}
	default:
		err := fmt.Errorf("unsupported action: %s", action)
		k.logger.Error(err)
		return err
	}

	if jsonData == nil {
		k.logger.Warn("No data found for action:", action)
		c.JSON(http.StatusOK, gin.H{
			"brand":  k.getBrand(),
			"id":     k.getID(),
			"action": action,
			"error":  err,
			"result": &jsonData,
		})
		return nil
	}

	k.logger.Debug("Action executed successfully:", action)
	c.JSON(200, gin.H{
		"brand":  k.getBrand(),
		"id":     k.getID(),
		"action": action,
		"result": jsonData,
		"status": "success",
	})
	return nil
}

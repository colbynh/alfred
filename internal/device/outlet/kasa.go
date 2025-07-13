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
	IPs []string `json:"ips"`
}

// Network scanning constants
const (
	timeout     = 1000 * time.Millisecond // Timeout for checking each port
	port1       = "9999"                  // Legacy Kasa device port
	port2       = "20002"                 // Newer Kasa device port
	subnet      = "192.168.101."
	IpBatchSize = 100 // Network subnet to scan
	startIP     = 1   // Focus on known device IPs
	endIP       = 254 // Focus on known device IPs
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

	ports := []string{port1, port2}

	for _, port := range ports {
		address := joinHostPort(ip, port)
		conn, err := dialTimeout("tcp", address, timeout)
		if err == nil && conn != nil {
			conn.Close()
			select {
			case results <- ip:
			default:
			}
			return
		}
	}
}

func ScanOpenPorts() ([]string, error) {
	startTime := time.Now()

	var wg sync.WaitGroup
	results := make(chan string, endIP-startIP+1)
	var openIPs []string

	for i := startIP; i <= endIP; i++ {
		ip := fmt.Sprintf("%s%d", subnet, i)
		wg.Add(1)
		go scanIP(ip, &wg, results)
		time.Sleep(100 * time.Millisecond)
	}

	// Create a channel to signal completion
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(results)
		close(done)
	}()

	// Wait for completion with timeout
	select {
	case <-done:
		// All goroutines completed
	case <-time.After(10 * time.Second):
		// Timeout - collect what we have so far
	}

	// Collect results
	for {
		select {
		case ip, ok := <-results:
			if !ok {
				goto scanComplete
			}
			openIPs = append(openIPs, ip)
		case <-time.After(2 * time.Second):
			// Additional timeout for collecting results
			goto scanComplete
		}
	}

scanComplete:
	elapsed := time.Since(startTime)
	fmt.Printf("Port scanning completed in %v\n", elapsed)

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

func (k *kasaOutlet) discoverDevicesKasa() (map[string]interface{}, error) {
	k.logger.Debug("Scanning for devices with Kasa tool on subnet:", subnet)
	startTime := time.Now()

	var wg sync.WaitGroup
	results := make(chan string, 100)
	start := 1
	end := 254

	for batchStart := start; batchStart <= end; batchStart += IpBatchSize {
		discoveryTimeout := "2"
		batchEnd := batchStart + IpBatchSize - 1
		if batchEnd > end {
			batchEnd = end
		}

		fmt.Printf("Scanning batch %d-%d\n", batchStart, batchEnd)

		for i := batchStart; i <= batchEnd; i++ {
			wg.Add(1)
			go func(ipNum int) {
				defer wg.Done()

				ip := fmt.Sprintf("%s%d", subnet, ipNum)
				cmd := execCommand("kasa", "--host", ip, "--discovery-timeout", discoveryTimeout, "state")
				o, err := cmd.Output()
				if err != nil {
					fmt.Printf("Error scanning %s: %v\n\n", ip, err)
					return
				}

				re := regexp.MustCompile(`Device state:\s+(False|True)`)
				match := re.FindString(string(o))

				if match != "" {
					fmt.Println("Found device:", ip)
					results <- ip
				}
			}(i)
		}
		wg.Wait()
		fmt.Printf("Completed batch %d-%d\n", batchStart, batchEnd)
	}

	close(results)

	var foundDevices []string
	for ip := range results {
		foundDevices = append(foundDevices, ip)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Kasa discovery completed in %v\n", elapsed)
	fmt.Printf("Found %d devices: %v\n", len(foundDevices), foundDevices)

	sr := ScanResult{IPs: foundDevices}
	jsonBytes, err := json.Marshal(sr)
	if err != nil {
		fmt.Printf("Error marshaling ScanResult: %v\n", err)
		return nil, err
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return nil, err
	}

	fmt.Printf("Discovered with kasas devices JSON output: %s\n", string(jsonBytes))

	return jsonData, nil
}

// discoverDevicesIps scans the network for Kasa devices and returns their IP addresses.
// The result is formatted as a JSON object with an "ips" array.
func (k *kasaOutlet) discoverDevicesIps() (map[string]interface{}, error) {
	k.logger.Debug("Scanning for open ports on subnet:", subnet)
	var ips []string
	var err error

	for attempts := 0; attempts < 3; attempts++ {
		k.logger.Debugf("Starting scan attempt %d", attempts+1)
		ips, err = ScanOpenPorts()
		if err == nil && len(ips) > 0 {
			k.logger.Debugf("Scan attempt %d successful, found %d devices", attempts+1, len(ips))
			break
		}
		if attempts < 2 {
			k.logger.Debugf("Scan attempt %d failed, retrying in 5 seconds", attempts+1)
			time.Sleep(time.Second * 1)
		}
	}

	if err != nil {
		k.logger.Warnf("Scan completed with error: %v", err)
	} else {
		k.logger.Debugf("Scan completed successfully, found %d devices", len(ips))
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

		for attempt := 1; attempt <= 3; attempt++ {
			cmd := execCommand("kasa", "--host", k.id, "--timeout", "10", "on")
			output, err := cmd.CombinedOutput()
			if err == nil {
				k.logger.Debugf("Kasa on command output: %s", string(output))
				break
			}
			k.logger.Warnf("Attempt %d failed: %v\nOutput: %s", attempt, err, string(output))
			if attempt < 3 {
				time.Sleep(time.Second * time.Duration(attempt))
			} else {
				k.logger.Errorf("All attempts failed for kasa on command: %v\nOutput: %s", err, string(output))
				return err
			}
		}
	case "off":
		k.logger.Debug("Turning off the device")

		for attempt := 1; attempt <= 3; attempt++ {
			cmd := execCommand("kasa", "--host", k.id, "--timeout", "10", "off")
			output, err := cmd.CombinedOutput()
			if err == nil {
				k.logger.Debugf("Kasa off command output: %s", string(output))
				break
			}
			k.logger.Warnf("Attempt %d failed: %v\nOutput: %s", attempt, err, string(output))
			if attempt < 3 {
				time.Sleep(time.Second * time.Duration(attempt))
			} else {
				k.logger.Errorf("All attempts failed for kasa off command: %v\nOutput: %s", err, string(output))
				return err
			}
		}
	case "discoverByKasa":
		k.logger.Debug("Discovering devices using kasa tool...")
		jsonData, err = k.discoverDevicesKasa()
		if err != nil {
			k.logger.Error("Error discovering devices:", err)
			return err
		}
	case "discoverByPorts":
		k.logger.Debug("Discovering devices using port scanning...")
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

	if c != nil {
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
	}
	return nil
}

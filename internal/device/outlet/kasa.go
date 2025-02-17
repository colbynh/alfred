package outlet

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	logrus "github.com/sirupsen/logrus"
)

type kasaOutlet struct {
	id     string
	c      *gin.Context
	logger *logrus.Logger
}

type ScanResult struct {
	IPs []string `json:"ips"`
}

func (k *kasaOutlet) getID() string {
	return k.id
}

func (k *kasaOutlet) getBrand() string {
	return "kasa"
}

func (k *kasaOutlet) discoverDevicesIps() (map[string]interface{}, error) {
	k.logger.Debug("Executing nmap command to discover devices")
	cmd := exec.Command("sh", "-c", `nmap -p 9999 --open --min-rate 10 100 192.168.101.0/24 | grep "Nmap scan report" | awk '{print $5}'`)

	output, err := cmd.Output()
	k.logger.Debug("nmap command string:", string(cmd.String()))

	if err != nil {
		k.logger.Error("Error executing nmap:", err)
		return nil, err
	}

	result := strings.TrimSpace(string(output))
	k.logger.Debug("nmap command output:", result)

	var ipList []string
	if result != "" {
		ipList = strings.Split(result, "\n")
	} else {
		k.logger.Warn("No devices found")
		return nil, errors.New("no devices found")
	}

	sr := ScanResult{IPs: ipList}
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

func (k *kasaOutlet) state() (map[string]interface{}, error) {
	k.logger.Debug("Executing kasa state command")
	cmd := exec.Command("kasa", "--host", k.id, "state")

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

func (k *kasaOutlet) sysInfo() (map[string]interface{}, error) {
	k.logger.Debug("Executing kasa sysinfo command")
	cmd := exec.Command("kasa", "--host", k.id, "sysinfo")

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

func (k *kasaOutlet) action(action string, c *gin.Context) error {
	k.logger.Debug("Executing action:", action)
	jsonData := map[string]interface{}{}
	var err error

	switch action {
	case "on":
		k.logger.Debug("Turning on the device")
		cmd := exec.Command("kasa", "--host", k.id, "on")

		if err = cmd.Run(); err != nil {
			k.logger.Error("Error executing kasa on command:", err)
			return err
		}
	case "off":
		k.logger.Debug("Turning off the device")
		cmd := exec.Command("kasa", "--host", k.id, "off")
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
	})
	return nil
}

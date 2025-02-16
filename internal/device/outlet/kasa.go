package outlet

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

type kasaOutlet struct {
	id string
	c  *gin.Context
}

func (k *kasaOutlet) getID() string {
	return k.id
}

func (k *kasaOutlet) getBrand() string {
	return "kasa"
}

func (k *kasaOutlet) state() (map[string]interface{}, error) {
	cmd := exec.Command("kasa", "--host", k.id, "state")

	o, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}

	re := regexp.MustCompile(`Device state:\s+(False|True)`) // Matches everything inside curly braces
	match := re.FindString(string(o))

	if match != "" {
		match = strings.ReplaceAll(match, "Device state:", "")
		jsonStr := fmt.Sprintf("{\"state\": \"%s\"}", strings.Trim(match, " "))

		if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
			return nil, err
		}
	}
	return jsonData, nil
}

func (k *kasaOutlet) sysInfo() (map[string]interface{}, error) {
	cmd := exec.Command("kasa", "--host", k.id, "sysinfo")

	o, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var jsonStr string
	var jsonData map[string]interface{}

	re := regexp.MustCompile(`(?s)\{.*\}`) // Matches everything inside curly braces
	match := re.FindString(string(o))

	if match != "" {
		// Replace single quotes with double quotes for valid JSON
		jsonStr = strings.ReplaceAll(match, "'", "\"")

		// Ensure nested JSON structures are properly formatted
		jsonStr = strings.ReplaceAll(jsonStr, "\"{", "{")
		jsonStr = strings.ReplaceAll(jsonStr, "}\"", "}")
		if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
			return nil, err
		}
	}

	return jsonData, nil
}

func (k *kasaOutlet) action(action string, c *gin.Context) error {
	jsonData := map[string]interface{}{}
	err := error(nil)
	switch action {
	case "on":
		cmd := exec.Command("kasa", "--host", k.id, "on")
		if err = cmd.Run(); err != nil {
			return err
		}
	case "off":
		cmd := exec.Command("kasa", "--host", k.id, "off")
		if err = cmd.Run(); err != nil {
			return err
		}
	case "state":
		jsonData, err = k.state()
		if err != nil {
			return err
		}
	case "sysinfo":
		jsonData, err = k.sysInfo()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported action: %s", action)
	}

	if action == "state" || action == "sysinfo" {
		if err != nil {
			fmt.Println("Error binding json", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"brand":  "kasa",
			"id":     "1234",
			"action": action,
			"error":  err,
			"result": &jsonData,
		})
		return nil
	}
	c.JSON(200, gin.H{
		"brand":  k.getBrand(),
		"id":     k.getID(),
		"action": action,
		"error":  err,
	})
	return nil
}

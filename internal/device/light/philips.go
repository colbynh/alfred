package light

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type philipsLight struct {
	brand      string
	ip         string
	id         string
	actionName string
	ctx        *gin.Context
}

func (p *philipsLight) getIP() string {
	return p.ip
}

func (p *philipsLight) getID() string {
	return p.id
}

func (p *philipsLight) getBrand() string {
	return "philips"
}

func (p *philipsLight) getGinContext() *gin.Context {
	return p.ctx
}

func (p *philipsLight) execAction(action string) error {
	switch action {
	case "getAll":
		p.actionName = "getAll"
		err := p.getAll()
		if err != nil {
			return err
		}
	case "on":
		p.actionName = "on"
		err := p.on()
		if err != nil {
			return err
		}
	case "off":
		p.actionName = "off"
		err := p.off()
		if err != nil {
			return err
		}
	case "brightness":
		p.actionName = "brightness"
		err := p.setBrightness()
		if err != nil {
			return err
		}
	case "color":
		p.actionName = "color"
		err := p.color()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported action: %s", action)
	}
	return nil
}

func (p *philipsLight) getAll() error {
	return runGetRequest(p)
}

func (p *philipsLight) on() error {
	jsonString := `{"on":{"on":true}}`

	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &jsonData); err != nil {
		fmt.Printf("unmarshal error: %s\n", err)
		return err
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Printf("marshal error: %s\n", err)
		return err
	}

	p.ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
	return runPutRequest(p)
}

func (p *philipsLight) off() error {
	jsonString := `{"on":{"on":false}}`

	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &jsonData); err != nil {
		fmt.Printf("unmarshal error: %s\n", err)
		return err
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Printf("marshal error: %s\n", err)
		return err
	}

	p.ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
	return runPutRequest(p)
}

func (p *philipsLight) setBrightness() error {
	// runRequest(p)
	return nil
}

func (p *philipsLight) color() error {
	return nil
}

// Helpers

func runPutRequest(p *philipsLight) error {
	url := fmt.Sprintf("https://%s/clip/v2/resource/light/%s", p.ip, p.id)

	bodyBytes, err := io.ReadAll(p.ctx.Request.Body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("hue-application-key", p.ctx.GetHeader("hue-application-key"))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to perform action: %s", p.actionName)
	}
	return nil
}

func runPostRequest(p *philipsLight) error {
	url := fmt.Sprintf("https://%s/clip/v2/resource/light/%s", p.ip, p.id)

	bodyBytes, err := io.ReadAll(p.ctx.Request.Body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("hue-application-key", p.ctx.GetHeader("hue-application-key"))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to perform action: %s", p.actionName)
	}
	return nil
}

func runGetRequest(p *philipsLight) error {
	var url string
	if p.actionName == "getAll" {
		url = fmt.Sprintf("https://%s/clip/v2/resource/light", p.ip)
	} else {
		url = fmt.Sprintf("https://%s/clip/v2/resource/light/%s", p.ip, p.id)
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("hue-application-key", p.ctx.GetHeader("hue-application-key"))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to perform action: %s", p.actionName)
	}
	return nil
}

func runDeleteRequest(p *philipsLight) error {
	var jsonData map[string]interface{}
	if err := p.ctx.ShouldBindJSON(&jsonData); err != nil {
		p.ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return err
	}

	url := fmt.Sprintf("https://%s/clip/v2/resource/light/%s", p.ip, p.id)

	jsonBytes, err := json.Marshal(p.ctx.Request.Body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("hue-application-key", p.ctx.GetHeader("hue-application-key"))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to perform action: %s", p.actionName)
	}
	return nil
}

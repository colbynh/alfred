package light

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
)

type philipsLight struct {
	brand       string
	ip          string
	id          string
	jsonPayload string
	actionName  string
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

func (p *philipsLight) execAction(action string) error {
	switch action {
	case "on":
		p.actionName = "on"
		p.on()
	case "off":
		p.actionName = "off"
		p.off()
	case "brightness":
		p.actionName = "brightness"
		p.setBrightness()
	case "color":
		p.actionName = "color"
		p.color()
	default:
		return fmt.Errorf("unsupported action: %s", action)
	}
	return nil
}

func (p *philipsLight) on() error {
	p.jsonPayload = `{"on":{"on":true}}`

	runRequest(p)
	return nil
}

func (p *philipsLight) off() error {
	p.jsonPayload = `"on":{"on":false}`
	runRequest(p)
	return nil
}

func (p *philipsLight) setBrightness() error {
	// runRequest(p)
	return nil
}

func (p *philipsLight) color() error {
	return nil
}

// Helpers

func runRequest(p *philipsLight) error {
	url := fmt.Sprintf("https://%s/clip/v2/resource/light/%s", p.ip, p.id)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(p.jsonPayload)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("hue-application-key", "<my-key>")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	requestDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return err
	}
	fmt.Printf("*****Request insecure: %s\n RequestBody: %s\n", string(requestDump), req.Body)

	resp, err := client.Do(req)

	respDump, _ := httputil.DumpResponse(resp, true)

	fmt.Printf("*****Responsedump: %s\n", respDump)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to perform action: %s", p.actionName)
	}
	return nil
}

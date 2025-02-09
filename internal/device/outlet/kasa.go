package outlet

import (
	"fmt"
	"os/exec"
	"strconv"
)

type kasaOutlet struct {
	id int
}

// GetID returns the outlet ID
func (k *kasaOutlet) getID() int {
	return k.id
}

// GetBrand returns the brand name
func (k *kasaOutlet) getBrand() string {
	return "kasa"
}

func (k *kasaOutlet) action(action string) error {
	id := strconv.Itoa(k.id)
	cmd := exec.Command("kasa", "--host", id, action)
	output, err := cmd.Output()

	fmt.Println(output)

	if err != nil {
		return err
	}
	return nil
}

package outlet

import (
	"fmt"
	"os/exec"
)

type kasaOutlet struct {
	id string
}

func (k *kasaOutlet) getID() string {
	return k.id
}

func (k *kasaOutlet) getBrand() string {
	return "kasa"
}

func (k *kasaOutlet) action(action string) error {
	cmd := exec.Command("kasa", "--host", k.id, action)

	_, err := cmd.Output()

	if err != nil {
		fmt.Printf("****cmd: %s\n\n *****err: %s\n", cmd, err)
		return err
	}
	return nil
}

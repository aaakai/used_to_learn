package test

import (
	"fmt"
	"os/exec"
	"strings"
)

func text() {
	cmd := exec.Command("ls", "-lh")
	std, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(std))

	fmt.Println(strings.Replace(string(std), "\n", "", -1))
}

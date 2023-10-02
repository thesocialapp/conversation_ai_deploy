package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("python", "main.py")
	cmd.Dir = "/app/internal/py/speech_recognition"
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	println(string(output))
	fmt.Printf("The type of output is %T\n", output)
}

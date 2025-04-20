// Utility functions
package internal

import (
	"fmt"
	"os"
	"bufio"
)

func ShowLogo() {
	file, err := os.Open("assets/logo.txt")
	if err != nil {
		fmt.Println("Welcome to CybeRPG!")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

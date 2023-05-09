package Token

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func OpenFile() string {
	configFile, err := os.Open("H:\\VKbot\\telegram\\Token\\token")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()
	input := bufio.NewScanner(configFile)
	input.Scan()
	parts := strings.Split(input.Text(), "=")
	token := strings.TrimSpace(parts[1])
	return token
}

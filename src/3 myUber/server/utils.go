package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

func listenOnPort() (lis net.Listener, port int, err error) {
	for {
		port = rand.Intn(65535-1024) + 1024
		lis, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			break
		}

		log.Printf("failed to listen on port %d: %v", port, err)
	}

	return lis, port, nil
}

func appendPortToFile(port int, portFilePath string) error {
	file, err := os.OpenFile(portFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%d\n", port))
	return err
}

func removePortFromFile(port int, portFilePath string) error {
	content, err := os.ReadFile(portFilePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != strconv.Itoa(port) {
			newLines = append(newLines, line)
		}
	}

	newContent := strings.Join(newLines, "\n")
	return os.WriteFile(portFilePath, []byte(newContent), 0644)
}
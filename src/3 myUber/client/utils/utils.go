package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func ReadPortsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open port file: %v", err)
	}
	defer file.Close()

	var ports []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ports = append(ports, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading port file: %v", err)
	}

	// get counts of each port
	portCount := make(map[string]int)
	for _, port := range ports {
		if len(port) == 0 {
			continue
		}
		portCount[port]++
	}

	// remove duplicates
	ports = nil
	for port, count := range portCount {
		if count > 1 {
			log.Printf("[warning] port %s found %d times", port, count)
		}
		ports = append(ports, port)
	}

	return ports, nil
}

package utils

import (
	"strconv"
	"strings"
)

func ParsePorts(portFlagInput string) (ports []int, err error) {
	splitInput := strings.Split(portFlagInput, ",")
	for _, part := range splitInput {
		if strings.Contains(part, "-") {
			portRange, err := getRange(part)
			if err != nil {
				return nil, err
			}

			ports = append(ports, portRange...)
			continue
		}

		port, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		ports = append(ports, port)
	}
	return
}

func getRange(portRange string) (ports []int, err error) {
	endpoints := strings.Split(portRange, "-")

	start, err := strconv.Atoi(endpoints[0])
	if err != nil {
		return nil, err
	}
	end, err := strconv.Atoi(endpoints[1])
	if err != nil {
		return nil, err
	}

	for i := start; i <= end; i++ {
		ports = append(ports, i)
	}
	return
}

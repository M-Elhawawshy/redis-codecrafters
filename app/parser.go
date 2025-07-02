package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func parseRESP(command string) ([]string, error) {
	switch command {
	case "PING":
		return []string{"PING"}, nil
	default:
		reader := bufio.NewReader(strings.NewReader(command))
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(line, "*") {
			return nil, fmt.Errorf("invalid RESP command")
		}
		nElms, err := strconv.Atoi(strings.TrimSpace(line[1:]))
		if err != nil {
			return nil, err
		}
		var res []string
		for _ = range nElms {
			line, err := reader.ReadString('\n')
			if err != nil {
				return nil, err
			}
			if !strings.HasPrefix(line, "$") {
				return nil, fmt.Errorf("expected bulk string, got: %s", line)
			}
			length, err := strconv.Atoi(strings.TrimSpace(line[1:]))
			if err != nil {
				return nil, err
			}

			buf := make([]byte, length+2) // +2 for \r\n
			_, err = io.ReadFull(reader, buf)
			if err != nil {
				return nil, err
			}
			str := string(buf[:length])
			res = append(res, str)
		}
		return res, nil
	}
}

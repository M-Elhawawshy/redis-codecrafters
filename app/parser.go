package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
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

func parseDatabaseRDB(reader *bufio.Reader) (map[string]ValueWithExpiry, error) {
	// get to db section
	for {
		b, err := reader.ReadByte()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if b == 0xFE {
			break
		}
	}
	// db index
	_, err := readLength(reader)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	hashTable, err := readLength(reader)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if hashTable != 0xFB {
		fmt.Println("invalid hash section")
		return nil, nil
	}
	hashSize, _ := readLength(reader)
	expirySize, _ := readLength(reader)
	fmt.Println("Hash size:", hashSize, " Expiring keys:", expirySize)
	keyValueMap := make(map[string]ValueWithExpiry)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		if b == 0xFF {
			fmt.Println("Reached EOF marker")
			break
		}

		var expiresAt time.Time

		if b == 0xFC || b == 0xFD {
			if b == 0xFC {
				ts := make([]byte, 8)
				io.ReadFull(reader, ts)
				ms := binary.LittleEndian.Uint64(ts)
				expiresAt = time.UnixMilli(int64(ms)) // Convert millis to time.Time
			} else {
				ts := make([]byte, 4)
				io.ReadFull(reader, ts)
				secs := binary.LittleEndian.Uint32(ts)
				expiresAt = time.Unix(int64(secs), 0) // Convert seconds to time.Time
			}
			b, _ = reader.ReadByte() // move on to value type
		}

		if b == 0x00 {
			key, _ := readString(reader)
			value, _ := readString(reader)
			keyValueMap[key] = ValueWithExpiry{
				Value:     value,
				ExpiresAt: expiresAt,
			}
		} else {
			fmt.Println("Unsupported type:", b)
		}
	}
	return keyValueMap, nil
}

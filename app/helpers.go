package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func readLength(r *bufio.Reader) (int, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	prefix := b >> 6

	switch prefix {
	case 0b00:
		return int(b & 0x3F), nil
	case 0b01:
		b2, _ := r.ReadByte()
		return int(b&0x3F)<<8 | int(b2), nil
	case 0b10:
		buf := make([]byte, 4)
		io.ReadFull(r, buf)
		return int(binary.BigEndian.Uint32(buf)), nil
	default:
		return 0, fmt.Errorf("unsupported length prefix: %02x", b)
	}
}

func readString(r *bufio.Reader) (string, error) {
	length, err := readLength(r)
	if err != nil {
		return "", err
	}
	buf := make([]byte, length)
	_, err = io.ReadFull(r, buf)
	return string(buf), err
}

func loadRDBToMemory() {
	app.config.RLock()
	dir := app.config.settings["dir"]
	filename := app.config.settings["dbfilename"]
	app.config.RUnlock()
	if dir != "" && filename != "" {
		file, err := os.Open(fmt.Sprintf("%s/%s", dir, filename))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		mp, err := parseDatabaseRDB(reader)
		if err != nil {
			fmt.Println("failed to parse db: ", err)
			return
		}

		app.db.Lock()
		for k, v := range mp {
			app.db.memory[k] = v
		}
		app.db.Unlock()
	}
}

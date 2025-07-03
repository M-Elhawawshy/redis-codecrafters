package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
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

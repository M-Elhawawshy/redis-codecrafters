package main

import (
	"bufio"
	"fmt"
	flag "github.com/spf13/pflag"
	"net"
	"os"
	"sync"
	"time"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

type DB struct {
	memory map[string]ValueWithExpiry
	sync.RWMutex
}

type ValueWithExpiry struct {
	Value     string
	ExpiresAt time.Time
}

type ConfigDB struct {
	sync.RWMutex
	settings map[string]string
}

type App struct {
	db     DB
	config ConfigDB
}

var app = App{
	db: DB{
		memory:  make(map[string]ValueWithExpiry),
		RWMutex: sync.RWMutex{},
	},
	config: ConfigDB{
		RWMutex:  sync.RWMutex{},
		settings: make(map[string]string),
	},
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	dir := flag.StringP("dir", "d", "", "specify the dir of persistent DB")
	dbfilename := flag.StringP("dbfilename", "", "", "specify the dir of persistent DB")
	flag.Parse()
	app.config.Lock()
	app.config.settings["dir"] = *dir
	app.config.settings["dbfilename"] = *dbfilename
	app.config.Unlock()
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go processConn(conn)
	}
}

func processConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	buf := make([]byte, 1024)
	for {
		_, err := reader.Read(buf)
		if err == nil {
			command, err := parseRESP(string(buf))
			if err != nil {
				fmt.Println("Error parsing command: ", err.Error())
				os.Exit(1)
			}
			processCommand(command, conn)
		}
	}
}

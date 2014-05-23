package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	tag string
	mu  sync.Mutex
)

func init() {
	tag = os.Args[0]
}

func SetTag(t string) {
	tag = t
}

func Debug(msg string) {
	write("DEBUG", msg)
}

func Error(msg string) {
	write("ERROR", msg)
}

func Info(msg string) {
	write("INFO", msg)
}

func Notice(msg string) {
	write("NOTICE", msg)
}

func Warning(msg string) {
	write("WARNING", msg)
}

func write(level, msg string) {
	timestamp := time.Now().Format(time.RFC3339)
	hostname, _ := os.Hostname()
	mu.Lock()
	defer mu.Unlock()

	fmt.Fprintf(os.Stderr, "%s %s %s[%d]: %s %s\n",
		timestamp, hostname, tag, os.Getpid(), level, msg)
}

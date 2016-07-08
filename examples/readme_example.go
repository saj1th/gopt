package main

import (
	"os"
	"time"

	"github.com/saj1th/gopt"
)

func main() {
	options := struct {
		Server   string        `gopt:"-s, --server, obligatory, description='Server to connect to'"`
		Password string        `gopt:"-p, --password, description='Don\\'t prompt for password'"`
		Timeout  time.Duration `gopt:"-t, --timeout, description='Connection timeout in seconds'"`
		Help     gopt.Help     `gopt:"-h, --help, description='Show this help'"`

		gopt.Verbs
		Execute struct {
			Command string   `gopt:"--command, mutexgroup='input', description='Command to exectute', obligatory"`
			Script  *os.File `gopt:"--script, mutexgroup='input', description='Script to exectute', rdonly"`
		} `gopt:"execute"`
		Delete struct {
			Path  string `gopt:"-n, --name, obligatory, description='Name of the entity to be deleted'"`
			Force bool   `gopt:"-f, --force, description='Force removal'"`
		} `gopt:"delete"`
	}{ // Default values goes here
		Timeout: 10 * time.Second,
	}
	gopt.ParseAndFail(&options)
}

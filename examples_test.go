package gopt

import (
	"fmt"
	"os"
	"time"
)

func ExampleFlagSet_PrintHelp() {
	options := struct {
		Server   string        `gopt:"-s, --server, obligatory, description='Server to connect to'"`
		Password string        `gopt:"-p, --password, description='Don\\'t prompt for password'"`
		Timeout  time.Duration `gopt:"-t, --timeout, description='Connection timeout in seconds'"`
		Help     Help          `gopt:"-h, --help, description='Show this help'"`

		Verbs
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

	args := []string{"--help"}
	fs := NewFlagSet("gopt", &options)
	err := fs.Parse(args)
	if err == ErrHelpRequest {
		fs.PrintHelp(os.Stdout)
		return
	} else if err != nil {
		fmt.Printf("Failure: %s", err)
	}

	// Output:
	// Usage: gopt [global options] <verb> [verb options]
	//
	// Global options:
	//         -s, --server   Server to connect to (*)
	//         -p, --password Don't prompt for password
	//         -t, --timeout  Connection timeout in seconds (default: 10s)
	//         -h, --help     Show this help
	//
	// Verbs:
	//     delete:
	//         -n, --name     Name of the entity to be deleted (*)
	//         -f, --force    Force removal
	//     execute:
	//             --command  Command to exectute (*)
	//             --script   Script to exectute
}

func ExampleVerbs() {
	options := struct {
		ImportantFlag string        `gopt:"-f, --flag, description='Important flag, obligatory'"`
		Password      string        `gopt:"-p, --password, description='Don\\'t prompt for password'"`
		Timeout       time.Duration `gopt:"-t, --timeout, description='Connection timeout in seconds'"`
		Help          Help          `gopt:"-h, --help, description='Show this help'"`

		Verb    Verbs
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

	args := []string{"delete", "-n", "/usr/bin"}
	fs := NewFlagSet("gopt", &options)
	_ = fs.Parse(args)
	// Error handling omitted
	fmt.Printf("Selected verb: %s", options.Verb)

	// Output:
	// Selected verb: delete
}

func ExampleRemainder() {
	options := struct {
		Username  string `gopt:"-u, --user, obligatory, description='Name of the user'"`
		Remainder Remainder
	}{}

	args := []string{"-u", "surma", "some", "more", "args"}
	fs := NewFlagSet("gopt", &options)
	_ = fs.Parse(args)
	// Error handling omitted
	fmt.Printf("Remainder: %#v", options.Remainder)

	// Output:
	// Remainder: gopt.Remainder{"some", "more", "args"}
}

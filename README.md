`gopt`  is a fork of https://github.com/voxelbrain/goptions

# Example

```go
package main

import (
	"github.com/saj1th/gopt"
	"os"
	"time"
)

func main() {
	options := struct {
		Servers  []string      `gopt:"-s, --server, obligatory, description='Servers to connect to'"`
		Password string        `gopt:"-p, --password, description='Don\\'t prompt for password'"`
		Timeout  time.Duration `gopt:"-t, --timeout, description='Connection timeout in seconds'"`
		Help     gopt.Help `gopt:"-h, --help, description='Show this help'"`

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
```

```
$ go run examples/readme_example.go --help
Usage: a.out [global options] <verb> [verb options]

Global options:
        -s, --server   Servers to connect to (*)
        -p, --password Don't prompt for password
        -t, --timeout  Connection timeout in seconds (default: 10s)
        -h, --help     Show this help

Verbs:
    delete:
        -n, --name     Name of the entity to be deleted (*)
        -f, --force    Force removal
    execute:
            --command  Command to exectute (*)
            --script   Script to exectute
```

# Quick Reference

## gopt

Each field of your struct can be tagged with a `gopt`

```go
    FieldName type `gopt:"-S, --long, options..."`
```

Where the short options (`-S`) are declared with a single dash and
long options (`--long`) are declared with two dashes. Either or
both may be declared.

After the short/long option names are one or more of the following:

### Global Options

* description='...'
* obligatory
* mutexgroup='GROUP_NAME'

### os.File specific

* create
* append
* rdonly
* wronly
* rdwr
* excl
* sync
* trunc
* perm=0777

## Supported Types

* bool
* string
* float64
* float32
* int
* int64
* int32
* gopt.Help
* *os.File
* *net.TCPAddr
* *url.URL
* time.Duration



---
Version 2.5.11

package gopt

import (
	"fmt"
	"testing"
)

func TestParse_StringValue(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Name string `gopt:"--name, -n"`
	}
	expected := "SomeName"

	args = []string{"--name", "SomeName"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err != nil {
		t.Fatalf("Flag parsing failed: %s", err)
	}
	if options.Name != expected {
		t.Fatalf("Expected %s for options.Name, got %s", expected, options.Name)
	}

	options.Name = ""

	args = []string{"-n", "SomeName"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err != nil {
		t.Fatalf("Flag parsing failed: %s", err)
	}
	if options.Name != expected {
		t.Fatalf("Expected %s for options.Name, got %s", expected, options.Name)
	}
}

func TestParse_ObligatoryStringValue(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Name string `gopt:"-n, obligatory"`
	}
	args = []string{}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err == nil {
		t.Fatalf("Parsing should have failed.")
	}

	args = []string{"-n", "SomeName"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}

	expected := "SomeName"
	if options.Name != expected {
		t.Fatalf("Expected %s for options.Name, got %s", expected, options.Name)
	}
}

func TestParse_UnknownFlag(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Name string `gopt:"--name, -n"`
	}
	args = []string{"-k", "4"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err == nil {
		t.Fatalf("Parsing should have failed.")
	}
}

func TestParse_FlagCluster(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Fast    bool `gopt:"-f"`
		Silent  bool `gopt:"-q"`
		Serious bool `gopt:"-s"`
		Crazy   bool `gopt:"-c"`
		Verbose bool `gopt:"-v"`
	}
	args = []string{"-fqcv"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}

	if !(options.Fast &&
		options.Silent &&
		!options.Serious &&
		options.Crazy &&
		options.Verbose) {
		t.Fatalf("Unexpected value: %v", options)
	}
}

func TestParse_MutexGroup(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Create bool `gopt:"--create, mutexgroup='action'"`
		Delete bool `gopt:"--delete, mutexgroup='action'"`
	}
	args = []string{"--create", "--delete"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err == nil {
		t.Fatalf("Parsing should have failed.")
	}
}

func TestParse_HelpFlag(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Name string `gopt:"--name, -n"`
		Help `gopt:"--help, -h"`
	}
	args = []string{"-n", "SomeNone", "-h"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err != ErrHelpRequest {
		t.Fatalf("Expected ErrHelpRequest, got: %s", err)
	}

	args = []string{"-n", "SomeNone"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err != nil {
		t.Fatalf("Unexpected error returned: %s", err)
	}
}

func TestParse_Verbs(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Server string `gopt:"--server, -s"`

		Verbs
		Create struct {
			Name string `gopt:"--name, -n"`
		} `gopt:"create"`
	}

	args = []string{"-s", "127.0.0.1", "create", "-n", "SomeDocument"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}

	if !(options.Server == "127.0.0.1" &&
		options.Create.Name == "SomeDocument" &&
		options.Verbs == "create") {
		t.Fatalf("Unexpected value: %v", options)
	}
}

func TestParse_IntValue(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Limit int `gopt:"-l"`
	}

	args = []string{"-l", "123"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}

	if !(options.Limit == 123) {
		t.Fatalf("Unexpected value: %v", options)
	}
}

func TestParse_Remainder(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Limit int `gopt:"-l"`
		Remainder
	}

	args = []string{"-l", "123", "Something", "SomethingElse"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}

	if !(len(options.Remainder) == 2 &&
		options.Remainder[0] == "Something" &&
		options.Remainder[1] == "SomethingElse") {
		t.Fatalf("Unexpected value: %v", options)
	}
}

func TestParse_VerbRemainder(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Limit int `gopt:"-l"`
		Remainder

		Verbs
		Create struct {
			Fast bool `gopt:"-f"`
			Remainder
		} `gopt:"create"`
	}

	args = []string{"create", "-f", "Something", "SomethingElse"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}

	if !(len(options.Remainder) == 2 &&
		options.Remainder[0] == "Something" &&
		options.Remainder[1] == "SomethingElse" &&
		options.Verbs == "create") {
		t.Fatalf("Unexpected value: %v", options)
	}
}

func TestParse_NoRemainder(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Fast bool `gopt:"-f"`
	}

	args = []string{"-f", "Something", "SomethingElse"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err == nil {
		t.Fatalf("Parsing should have failed")
	}
}

func TestParse_MissingValue(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Name string `gopt:"-n, --name"`
	}

	args = []string{"-n"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err == nil {
		t.Fatalf("Parsing should have failed")
	}

	args = []string{"--name"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err == nil {
		t.Fatalf("Parsing should have failed")
	}
}

func TestParse_ObligatoryMutexGroup(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Create bool `gopt:"-c, mutexgroup='action', obligatory"`
		Delete bool `gopt:"-d, mutexgroup='action'"`
	}

	args = []string{}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err == nil {
		t.Fatalf("Parsing should have failed")
	}

	args = []string{"-c", "-d"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err == nil {
		t.Fatalf("Parsing should have failed")
	}

	args = []string{"-d"}
	fs = NewFlagSet("gopt", &options)
	err = fs.Parse(args)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
}

func TestParse_StringArray_Short(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Servers []string `gopt:"-s"`
	}

	args = []string{}
	for i := 1; i < 10; i++ {
		options.Servers = []string{}
		args = append(args, []string{"-s", fmt.Sprintf("server%d", i)}...)
		fs = NewFlagSet("gopt", &options)
		err = fs.Parse(args)
		if err != nil {
			t.Fatalf("Parsing failed at %d: %s", i, err)
		}
		if len(options.Servers) != i {
			t.Fatalf("Unexpected number of values. Expected %d, got %d (%#v)", i, len(options.Servers), options.Servers)
		}
		for j := 0; j < i; j++ {
			expected := fmt.Sprintf("server%d", j+1)
			if options.Servers[j] != expected {
				t.Fatalf("Unexpected value. %#v", options.Servers)
			}
		}
	}
}

func TestParse_BoolArray_Cluster(t *testing.T) {
	var err error
	var fs *FlagSet
	var options struct {
		Verbosity []bool `gopt:"-v"`
	}

	args := "-v"
	for i := 1; i < 10; i++ {
		options.Verbosity = []bool{}
		fs = NewFlagSet("gopt", &options)
		err = fs.Parse([]string{args})
		if err != nil {
			t.Fatalf("Parsing failed at %d: %s", i, err)
		}
		if len(options.Verbosity) != i {
			t.Fatalf("Unexpected number of values. Expected %d, got %d (%#v)", i, len(options.Verbosity), options.Verbosity)
		}
		args += "v"
	}
}

func TestParse_BoolArray_Short(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Verbosity []bool `gopt:"-v"`
	}

	args = []string{}
	for i := 1; i < 10; i++ {
		options.Verbosity = []bool{}
		args = append(args, "-v")
		fs = NewFlagSet("gopt", &options)
		err = fs.Parse(args)
		if err != nil {
			t.Fatalf("Parsing failed at %d: %s", i, err)
		}
		if len(options.Verbosity) != i {
			t.Fatalf("Unexpected number of values. Expected %d, got %d (%#v)", i, len(options.Verbosity), options.Verbosity)
		}
	}
}

func TestParse_BoolArray_Long(t *testing.T) {
	var args []string
	var err error
	var fs *FlagSet
	var options struct {
		Verbosity []bool `gopt:"--verbose"`
	}

	args = []string{}
	for i := 1; i < 10; i++ {
		options.Verbosity = []bool{}
		args = append(args, "--verbose")
		fs = NewFlagSet("gopt", &options)
		err = fs.Parse(args)
		if err != nil {
			t.Fatalf("Parsing failed at %d: %s", i, err)
		}
		if len(options.Verbosity) != i {
			t.Fatalf("Unexpected number of values. Expected %d, got %d (%#v)", i, len(options.Verbosity), options.Verbosity)
		}
	}
}

func TestParse_UnexportedVerbs(t *testing.T) {
	var options struct {
		Verbs
		A struct {
			A1 string `gopt:"--a1"`
			a2 string `gopt:"--a2"`
		} `gopt:"A"`
	}
	args := []string{"A", "--a1", "x"}
	fs := NewFlagSet("gopt", &options)
	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
	if options.A.A1 != "x" || options.A.a2 != "" {
		t.Fatalf("Unexpected values in struct: %#v", options)
	}
}

func TestParse_DashAsRemainder(t *testing.T) {
	var options struct {
		SomeFlag bool `gopt:"-b"`
		Remainder
	}
	args := []string{"-b", "-"}
	fs := NewFlagSet("gopt", &options)
	err := fs.Parse(args)
	if err != nil {
		t.Fatalf("Parsing failed: %s", err)
	}
	if len(options.Remainder) != 1 {
		t.Fatalf("Unexpected size of remainder: %d (%#v)", len(options.Remainder), options.Remainder)
	}
	if options.Remainder[0] != "-" {
		t.Fatalf("Unexpected remainder: %#v", options.Remainder)
	}

}

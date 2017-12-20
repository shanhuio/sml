package subcmd

import (
	"fmt"
	"io"
	"os"
	"sort"

	"shanhu.io/misc/httputil"
)

type command struct {
	name string
	desc string
	f    Func
	hf   HostFunc
}

// List is a list of sub command entires.
type List struct {
	cmds          map[string]*command
	defaultServer string
}

// New creates an empty list of sub command entries.
func New() *List {
	return &List{
		cmds: make(map[string]*command),
	}
}

// SetDefaultServer sets the default server for commands that has a host.
func (lst *List) SetDefaultServer(s string) {
	lst.defaultServer = s
}

// Add adds a sub command entry. It panics if already exists a command with the
// same name.
func (lst *List) Add(name, desc string, f Func) {
	if _, ok := lst.cmds[name]; ok {
		panic(fmt.Errorf("%q already exist", name))
	}

	c := &command{
		name: name,
		desc: desc,
		f:    f,
	}
	lst.cmds[name] = c
}

// AddHost adds a sub command entry that can accept a host target. It panics if
// already exists a command with the same name.
func (lst *List) AddHost(name, desc string, f HostFunc) {
	if _, ok := lst.cmds[name]; ok {
		panic(fmt.Errorf("%q already exist", name))
	}

	c := &command{
		name: name,
		desc: desc,
		hf:   f,
	}
	lst.cmds[name] = c
}

// Help prints out the help message.
func (lst *List) Help(w io.Writer) {
	var names []string
	for name := range lst.cmds {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		cmd := lst.cmds[name]
		fmt.Fprintf(w, "%s - %s\n", name, cmd.desc)
	}
}

// Main runs with os.Args, and calls os.Exit if the return value is not 0.
func (lst *List) Main() {
	if ret := lst.Run(os.Args); ret != 0 {
		os.Exit(ret)
	}
}

func (lst *List) parseServer(s string) string {
	if s == "" {
		return lst.defaultServer
	}
	return httputil.ExtendServer(s)
}

// Run runs and dispatches the sub command.
func (lst *List) Run(args []string) int {
	if len(args) < 2 {
		lst.Help(os.Stderr)
		return -1
	}

	name, host := parseCmd(args[1])
	if name == "-h" || name == "help" {
		lst.Help(os.Stdout)
		return 0
	}

	host = lst.parseServer(host)

	c, ok := lst.cmds[name]
	if !ok {
		fmt.Printf("command %q not found\n", name)
		return -1
	}

	if c.f != nil {
		if err := c.f(args[2:]); err != nil {
			fmt.Println(err)
			return -1
		}
	} else {
		if err := c.hf(host, args[2:]); err != nil {
			fmt.Println(err)
			return -1
		}
	}
	return 0
}

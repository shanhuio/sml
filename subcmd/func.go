package subcmd

// Func is an entry function of a sub command.
type Func func(args []string) error

// HostFunc is an entry function of a sub command with a host target.
type HostFunc func(host string, args []string) error

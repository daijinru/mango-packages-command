package mango_packages_command

import (
	"errors"
	"log"
	"os"
	"strings"
)

type PositionalArgs func(cmd *Command, args []string) error

func ExactArgs(n int) PositionalArgs {
	return func(cmd *Command, args []string) error {
		if len(args) != n {
			log.Fatalf("accepts %d arg(s), received %d", n, len(args))
		}
		return nil
	}
}

type Command struct {
	Use      string                                  `json:"use"`
	RunE     func(cmd *Command, args []string) error `json:"runE"`
	Args     PositionalArgs
	commands []*Command `json:"commands"`
	parent   *Command   `json:"parent"`
	args     []string   `json:"args"`
}

func (c *Command) ValidateArgs(args []string) error {
	if c.Args == nil {
		return nil
	}
	return c.Args(c, args)
}

func (c *Command) LogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Command) Execute() error {
	err := c.ExecuteC()
	return err
}

func (c *Command) FlagsString() string {
	var out string
	for _, f := range c.args {
		out = out + f + ","
	}
	return out
}

func (c *Command) ExecuteC() (err error) {
	args := os.Args[1:]
	cmd, flags, err := c.Find(args)
	if err != nil {
		c.LogFatal(err)
	}
	c.args = flags
	if cmd != nil {
		err = cmd.execute(flags)
		c.LogFatal(err)
	} else {
		c.LogFatal(errors.New(c.args[0] + " flag not exist: " + c.FlagsString()))
	}

	return err
}

func (c *Command) execute(a []string) error {
	if c.RunE == nil {
		c.LogFatal(errors.New("RunE does not exist"))
	}
	if err := c.ValidateArgs(a); err != nil {
		return err
	}
	err := c.RunE(c, a)
	c.LogFatal(err)
	return nil
}

func stripFlags(args []string) []string {
	if len(args) == 0 {
		return args
	}

	var commands []string

Loop:
	for len(args) > 0 {
		s := args[0]
		args = args[1:]
		switch {
		case s == "--":
			break Loop
		case strings.HasPrefix(s, "--"):
			fallthrough
		case strings.HasPrefix(s, "-"):
			if len(args) <= 1 {
				break Loop
			} else {
				args = args[1:]
				continue
			}
		case s != "" && !strings.HasPrefix(s, "-"):
			commands = append(commands, s)
		}
	}

	return commands
}

func (c *Command) Find(args []string) (*Command, []string, error) {
	// error: "Unresolved reference 'recFind'"
	// requires declaring a recursive func first
	var recFind func(*Command, []string) (*Command, []string)

	recFind = func(inC *Command, inArgs []string) (*Command, []string) {
		striped := stripFlags(inArgs)
		if len(striped) == 0 {
			return inC, inArgs
		}
		nextSubFlag := striped[0]
		cmd := findSub(nextSubFlag, inC)
		if cmd != nil {
			return recFind(cmd, striped[1:])
		}
		return inC, inArgs
	}

	commandFound, flags := recFind(c, args)
	return commandFound, flags, nil
}

func (c *Command) Name() string {
	name := c.Use
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func findSub(next string, parent *Command) *Command {
	for _, cmd := range parent.commands {
		if cmd.Name() == next {
			return cmd
		}
	}
	return nil
}

func (c *Command) AddCommand(cmd *Command) {
	cmd.parent = c
	c.commands = append(c.commands, cmd)
}

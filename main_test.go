package mango_packages_command

import "testing"

func TestMain(m *testing.M) {
	m.Run()
}

func TestSubCommand(t *testing.T) {
	root := &Command{
		Use: "RC",
		RunE: func(cmd *Command, args []string) error {
			return nil
		},
	}

	sub := &Command{
		Use:  "sub",
		Args: ExactArgs(1),
		RunE: func(cmd *Command, args []string) error {
			return nil
		},
	}

	root.AddCommand(sub)
}

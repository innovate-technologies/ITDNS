package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

type addOptions struct {
	Type string
}

// NewAddCmd generates the add command
func NewAddCmd() *cobra.Command {
	o := addOptions{}
	c := &cobra.Command{
		Use:     "add",
		Short:   "Adds/Replaces a DNS entry",
		RunE:    o.RunE,
		Example: `itdns add -t A shoutca.st 127.0.0.1`,
	}

	c.Flags().StringVarP(&o.Type, "type", "t", "A", "DNS record type")

	return c
}

func (a *addOptions) RunE(cmd *cobra.Command, args []string) error {
	if len(args) > 2 {
		return errors.New("Needs at least 2 argumants: name and value")
	}

	return nil
}

package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

type rmOptions struct {
	Type string
	All  bool
}

// NewRmCmd generates the add command
func NewRmCmd() *cobra.Command {
	o := rmOptions{}
	c := &cobra.Command{
		Use:   "rm",
		Short: "removes a DNS entry",
		RunE:  o.RunE,
		Example: `
## Remove only A record
itdns add -t A shoutca.st
## Remove ALL records for domain
itdns add --all shoutca.st
`,
	}

	c.Flags().StringVarP(&o.Type, "type", "t", "", "DNS record type")
	c.Flags().BoolVar(&o.All, "all", false, "Remove all records for domain")
	return c
}

func (a *rmOptions) RunE(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return errors.New("Needs at least 1 argumant: name")
	}

	return nil
}

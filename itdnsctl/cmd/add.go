package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type addOptions struct {
	Type   string
	TTL    int64
	DryRun bool
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
	c.Flags().Int64Var(&o.TTL, "ttl", 100, "TTL for the DNS record")
	c.Flags().BoolVar(&o.DryRun, "dry-run", false, "Shows the entry but will not create it")

	return c
}

func (a *addOptions) RunE(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("Needs at least 2 argumants: name and value")
	}

	// prepare content
	key := fmt.Sprintf("/DNS/%s/%s", args[0], strings.ToUpper(a.Type))
	value := fmt.Sprintf(`{"value": "%s", "ttl": %d}`, args[1], a.TTL)

	if a.DryRun {
		fmt.Println(key)
		fmt.Println(value)
		return nil
	}

	client, err := newClientFromEnv()
	if err != nil {
		return err
	}

	_, err = client.Put(context.Background(), key, value)
	return err
}

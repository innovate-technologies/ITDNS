package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/spf13/cobra"
)

type rmOptions struct {
	Type   string
	All    bool
	DryRun bool
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
	c.Flags().BoolVar(&o.DryRun, "dry-run", false, "Shows the entry but will not create it")

	return c
}

func (r *rmOptions) RunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("Needs at least 1 argumant: name")
	}

	if r.Type == "" && !r.All {
		return errors.New("Either --type or --all needs to be specified")
	}

	// prepare content
	var suffix string
	if !r.All {
		suffix = strings.ToUpper(r.Type)
	}
	key := fmt.Sprintf("/DNS/%s/%s", args[0], suffix)

	if r.DryRun {
		fmt.Println(key)
		return nil
	}

	client, err := newClientFromEnv()
	if err != nil {
		return err
	}

	_, err = client.Delete(context.Background(), key, etcd.WithPrefix())

	return nil
}

package cmd

import (
	"context"
	"encoding/json"
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

type record struct {
	Value string `json:"value"`
	TTL   int64  `json:"ttl"`
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

	value := []record{}

	// prepare content
	key := fmt.Sprintf("/DNS/%s/%s", args[0], strings.ToUpper(a.Type))
	value = append(value, record{
		Value: args[1],
		TTL:   a.TTL,
	})
	if len(args) > 2 {
		for _, val := range args[2:] {
			value = append(value, record{
				Value: val,
				TTL:   a.TTL,
			})
		}
	}

	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	valueString := string(valueBytes)

	if a.DryRun {
		fmt.Println(key)
		fmt.Println(valueString)
		return nil
	}

	client, err := newClientFromEnv()
	if err != nil {
		return err
	}

	_, err = client.Put(context.Background(), key, valueString)
	return err
}

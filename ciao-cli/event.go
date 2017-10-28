//
// Copyright (c) 2016 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ciao-project/ciao/ciao-controller/types"
	"github.com/intel/tfortools"
	"github.com/pkg/errors"
)

var eventCommand = &command{
	SubCommands: map[string]subCommand{
		"list":   new(eventListCommand),
		"delete": new(eventDeleteCommand),
	},
}

type eventListCommand struct {
	Flag     flag.FlagSet
	all      bool
	tenant   string
	template string
}

func (cmd *eventListCommand) usage(...string) {
	fmt.Fprintf(os.Stderr, `usage: ciao-cli [options] event list [flags]

List prints a list of events found in the ciao cluster

The list flags are:

`)
	cmd.Flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n%s",
		tfortools.GenerateUsageDecorated("f", types.CiaoEvents{}.Events, nil))
	os.Exit(2)
}

func (cmd *eventListCommand) parseArgs(args []string) []string {
	cmd.Flag.BoolVar(&cmd.all, "all", false, "List events for all tenants in a cluster")
	cmd.Flag.StringVar(&cmd.tenant, "tenant-id", "", "Tenant ID")
	cmd.Flag.StringVar(&cmd.template, "f", "", "Template used to format output")
	cmd.Flag.Usage = func() { cmd.usage() }
	cmd.Flag.Parse(args)
	return cmd.Flag.Args()
}

func (cmd *eventListCommand) run(args []string) error {
	if cmd.tenant == "" {
		cmd.tenant = c.TenantID
	}

	if cmd.all == false && cmd.tenant == "" {
		errorf("Missing required -tenant-id parameter")
		cmd.usage()
	}

	tenantID := cmd.tenant
	if cmd.all {
		tenantID = ""
	}

	events, err := c.ListEvents(tenantID)
	if err != nil {
		return errors.Wrap(err, "Error listing events")
	}

	if cmd.template != "" {
		return tfortools.OutputToTemplate(os.Stdout, "event-list", cmd.template,
			&events.Events, nil)
	}

	fmt.Printf("%d Ciao event(s):\n", len(events.Events))
	for i, event := range events.Events {
		fmt.Printf("\t[%d] %v: %s:%s (Tenant %s)\n", i+1, event.Timestamp, event.EventType, event.Message, event.TenantID)
	}
	return nil
}

type eventDeleteCommand struct {
	Flag flag.FlagSet
}

func (cmd *eventDeleteCommand) usage(...string) {
	fmt.Fprintf(os.Stderr, `Usage: ciao-cli [options] event delete

Deletes all events
`)
	os.Exit(2)
}

func (cmd *eventDeleteCommand) parseArgs(args []string) []string {
	cmd.Flag.Usage = func() { cmd.usage() }
	cmd.Flag.Parse(args)
	return cmd.Flag.Args()
}

func (cmd *eventDeleteCommand) run(args []string) error {
	err := c.DeleteEvents()
	if err != nil {
		return errors.Wrap(err, "Error deleting events")
	}
	fmt.Printf("Deleted all event logs\n")
	return nil
}

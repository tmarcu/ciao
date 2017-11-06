package sdk

import (
	"fmt"
	"os"

	"github.com/ciao-project/ciao/client"
	"github.com/intel/tfortools"
	"github.com/pkg/errors"
)

func ListEvents(c *client.Client, flags CommandOpts) error {
	if flags.Tenant == "" {
		flags.Tenant = c.TenantID
	}

	if flags.All == false && flags.Tenant == "" {
		client.Errorf("Missing required --tenantID parameter")
	}

	tenantID := flags.Tenant
	if flags.All {
		tenantID = ""
	}

	events, err := c.ListEvents(tenantID)
	if err != nil {
		return errors.Wrap(err, "Error listing events")
	}

	if c.Template != "" {
		return tfortools.OutputToTemplate(os.Stdout, "event-list", c.Template,
			&events.Events, nil)
	}

	fmt.Printf("%d Ciao event(s):\n", len(events.Events))
	for i, event := range events.Events {
		fmt.Printf("\t[%d] %v: %s:%s (Tenant %s)\n", i+1, event.Timestamp, event.EventType, event.Message, event.TenantID)
	}
	return nil
}

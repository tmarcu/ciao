package sdk

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ciao-project/ciao/client"
	"github.com/intel/tfortools"
	"github.com/pkg/errors"
)

func ListExternalIP(c *client.Client, flags CommandOpts) error {
	IPs, err := c.ListExternalIPs()
	if err != nil {
		return errors.Wrap(err, "Error listing external IPs")
	}

	if c.Template != "" {
		return tfortools.OutputToTemplate(os.Stdout, "external-ip-list", c.Template,
			&IPs, nil)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 1, 1, ' ', 0)
	fmt.Fprintf(w, "#\tExternalIP\tInternalIP\tInstanceID")
	if c.IsPrivileged() {
		fmt.Fprintf(w, "\tTenantID\tPoolName\n")
	} else {
		fmt.Fprintf(w, "\n")
	}

	for i, IP := range IPs {
		fmt.Fprintf(w, "%d", i+1)
		fmt.Fprintf(w, "\t%s", IP.ExternalIP)
		fmt.Fprintf(w, "\t%s", IP.InternalIP)
		if IP.InstanceID != "" {
			fmt.Fprintf(w, "\t%s", IP.InstanceID)
		}

		if IP.TenantID != "" {
			fmt.Fprintf(w, "\t%s", IP.TenantID)
		}

		if IP.PoolName != "" {
			fmt.Fprintf(w, "\t%s", IP.PoolName)
		}

		fmt.Fprintf(w, "\n")
	}

	w.Flush()

	return nil
}

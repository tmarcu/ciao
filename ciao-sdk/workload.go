package sdk

import (
	"fmt"
	"os"

	"github.com/ciao-project/ciao/client"
	"github.com/ciao-project/ciao/payloads"
	"github.com/intel/tfortools"
	"github.com/pkg/errors"

)

func ListWorkload(c *client.Client, flags CommandOpts) error {
	if c.TenantID == "" {
		fatalf("Missing required TenantID parameter")
	}

	wls, err := c.ListWorkloads()
	if err != nil {
		return errors.Wrap(err, "Error listing workloads")
	}

	var workloads []Workload
	for i, wl := range wls {
		workloads = append(workloads, Workload{
			Name: wl.Description,
			ID:   wl.ID,
		})

		for _, r := range wl.Defaults {
			if r.Type == payloads.MemMB {
				workloads[i].Mem = r.Value
			}
			if r.Type == payloads.VCPUs {
				workloads[i].CPUs = r.Value
			}
		}
	}

	if c.Template != "" {
		return tfortools.OutputToTemplate(os.Stdout, "workload-list", c.Template,
			workloads, nil)
	}

	for i, wl := range workloads {
		fmt.Printf("Workload %d\n", i+1)
		fmt.Printf("\tName: %s\n\tUUID:%s\n\tCPUs: %d\n\tMemory: %d MB\n",
			wl.Name, wl.ID, wl.CPUs, wl.Mem)
	}

	return nil
}
package sdk

import (
	"fmt"
	"os"

	"github.com/ciao-project/ciao/ciao-controller/types"
	"github.com/ciao-project/ciao/client"
	"github.com/ciao-project/ciao/payloads"

	"github.com/intel/tfortools"
	"github.com/pkg/errors"

	"gopkg.in/yaml.v2"
)

func outputWorkload(w types.Workload) {
	var opt workloadOptions

	opt.Description = w.Description
	opt.VMType = string(w.VMType)
	opt.FWType = w.FWType
	opt.ImageName = w.ImageName
	for _, d := range w.Defaults {
		if d.Type == payloads.VCPUs {
			opt.Defaults.VCPUs = d.Value
		} else if d.Type == payloads.MemMB {
			opt.Defaults.MemMB = d.Value
		}
	}

	for _, s := range w.Storage {
		d := disk{
			Size:      s.Size,
			Bootable:  s.Bootable,
			Ephemeral: s.Ephemeral,
		}
		if s.ID != "" {
			d.ID = &s.ID
		}

		src := source{
			Type: s.SourceType,
			ID:   s.SourceID,
		}

		d.Source = src

		opt.Disks = append(opt.Disks, d)
	}

	b, err := yaml.Marshal(opt)
	if err != nil {
		fatalf(err.Error())
	}

	fmt.Println(string(b))
	fmt.Println(w.Config)
}

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

func ShowWorkload(c *client.Client, flags CommandOpts) error {
	var wl types.Workload

	if len(flags.Args) == 0 {
		fatalf("Missing required workload UUID parameter")
	}
	workload := flags.Args[0]

	wl, err := c.GetWorkload(workload)
	if err != nil {
		return errors.Wrap(err, "Error getting workload")
	}

	if c.Template != "" {
		return tfortools.OutputToTemplate(os.Stdout, "workload-show", c.Template, &wl, nil)
	}

	outputWorkload(wl)
	return nil
}

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
	"io/ioutil"
	"os"

	"github.com/ciao-project/ciao/ciao-controller/types"
	"github.com/ciao-project/ciao/payloads"
	"github.com/intel/tfortools"
	"github.com/pkg/errors"

	"gopkg.in/yaml.v2"
)

var workloadCommand = &command{
	SubCommands: map[string]subCommand{
		"list":   new(workloadListCommand),
		"create": new(workloadCreateCommand),
		"delete": new(workloadDeleteCommand),
		"show":   new(workloadShowCommand),
	},
}

type workloadListCommand struct {
	Flag     flag.FlagSet
	template string
}

// Workload contains detailed information about a workload
type Workload struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	CPUs int    `json:"vcpus"`
	Mem  int    `json:"ram"`
}

func (cmd *workloadListCommand) usage(...string) {
	fmt.Fprintf(os.Stderr, `usage: ciao-cli [options] workload list

List all workloads

`)
	cmd.Flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n%s",
		tfortools.GenerateUsageDecorated("f", []Workload{}, nil))
	os.Exit(2)
}

func (cmd *workloadListCommand) parseArgs(args []string) []string {
	cmd.Flag.StringVar(&cmd.template, "f", "", "Template used to format output")
	cmd.Flag.Usage = func() { cmd.usage() }
	cmd.Flag.Parse(args)
	return cmd.Flag.Args()
}

func (cmd *workloadListCommand) run(args []string) error {
	if c.TenantID == "" {
		fatalf("Missing required -tenant-id parameter")
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

	if cmd.template != "" {
		return tfortools.OutputToTemplate(os.Stdout, "workload-list", cmd.template,
			workloads, nil)
	}

	for i, wl := range workloads {
		fmt.Printf("Workload %d\n", i+1)
		fmt.Printf("\tName: %s\n\tUUID:%s\n\tCPUs: %d\n\tMemory: %d MB\n",
			wl.Name, wl.ID, wl.CPUs, wl.Mem)
	}

	return nil
}

type workloadCreateCommand struct {
	Flag     flag.FlagSet
	yamlFile string
}

func (cmd *workloadCreateCommand) parseArgs(args []string) []string {
	cmd.Flag.StringVar(&cmd.yamlFile, "yaml", "", "filename for yaml which describes the workload")
	cmd.Flag.Usage = func() { cmd.usage() }
	cmd.Flag.Parse(args)
	return cmd.Flag.Args()
}

func (cmd *workloadCreateCommand) usage(...string) {
	fmt.Fprintf(os.Stderr, `usage: ciao-cli [options] workload create [flags]

Create a new workload

The create flags are:

`)
	cmd.Flag.PrintDefaults()
	os.Exit(2)
}

type source struct {
	Type types.SourceType `yaml:"service"`
	ID   string           `yaml:"id"`
}

type disk struct {
	ID        *string `yaml:"volume_id,omitempty"`
	Size      int     `yaml:"size"`
	Bootable  bool    `yaml:"bootable"`
	Source    source  `yaml:"source"`
	Ephemeral bool    `yaml:"ephemeral"`
}

type defaultResources struct {
	VCPUs int `yaml:"vcpus"`
	MemMB int `yaml:"mem_mb"`
}

// we currently only use the first disk due to lack of support
// in types.Workload for multiple storage resources.
type workloadOptions struct {
	Description     string           `yaml:"description"`
	VMType          string           `yaml:"vm_type"`
	FWType          string           `yaml:"fw_type,omitempty"`
	ImageName       string           `yaml:"image_name,omitempty"`
	Defaults        defaultResources `yaml:"defaults"`
	CloudConfigFile string           `yaml:"cloud_init,omitempty"`
	Disks           []disk           `yaml:"disks,omitempty"`
}

func optToReqStorage(opt workloadOptions) ([]types.StorageResource, error) {
	storage := make([]types.StorageResource, 0)
	bootableCount := 0
	for _, disk := range opt.Disks {
		res := types.StorageResource{
			Size:      disk.Size,
			Bootable:  disk.Bootable,
			Ephemeral: disk.Ephemeral,
		}

		// Use existing volume
		if disk.ID != nil {
			res.ID = *disk.ID
		} else {
			// Create a new one
			if disk.Source.Type == "" {
				disk.Source.Type = types.Empty
			}

			if disk.Source.Type != types.Empty {
				res.SourceType = disk.Source.Type
				res.SourceID = disk.Source.ID

				if res.SourceID == "" {
					return nil, errors.New("Invalid workload yaml: when using a source an id must also be specified")
				}
			} else {
				if disk.Bootable == true {
					// you may not request a bootable drive
					// from an empty source
					return nil, errors.New("Invalid workload yaml: empty disk source may not be bootable")
				}

				if disk.Size <= 0 {
					return nil, errors.New("Invalid workload yaml: size required when creating a volume")
				}
			}
		}

		if disk.Bootable {
			bootableCount++
		}

		storage = append(storage, res)
	}

	if payloads.Hypervisor(opt.VMType) == payloads.QEMU && bootableCount == 0 {
		return nil, errors.New("Invalid workload yaml: no bootable disks specified for a VM")
	}

	return storage, nil
}

func optToReq(opt workloadOptions, req *types.Workload) error {
	b, err := ioutil.ReadFile(opt.CloudConfigFile)
	if err != nil {
		return err
	}

	config := string(b)

	// this is where you'd validate that the options make
	// sense.
	req.Description = opt.Description
	req.VMType = payloads.Hypervisor(opt.VMType)
	req.FWType = opt.FWType
	req.ImageName = opt.ImageName
	req.Config = config
	req.Storage, err = optToReqStorage(opt)

	if err != nil {
		return err
	}

	// all default resources are required.
	defaults := opt.Defaults

	r := payloads.RequestedResource{
		Type:  payloads.VCPUs,
		Value: defaults.VCPUs,
	}
	req.Defaults = append(req.Defaults, r)

	r = payloads.RequestedResource{
		Type:  payloads.MemMB,
		Value: defaults.MemMB,
	}
	req.Defaults = append(req.Defaults, r)

	return nil
}

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

func (cmd *workloadCreateCommand) run(args []string) error {
	var opt workloadOptions
	var req types.Workload

	if cmd.yamlFile == "" {
		cmd.usage()
	}

	f, err := ioutil.ReadFile(cmd.yamlFile)
	if err != nil {
		fatalf("Unable to read workload config file: %s\n", err)
	}

	err = yaml.Unmarshal(f, &opt)
	if err != nil {
		fatalf("Config file invalid: %s\n", err)
	}

	err = optToReq(opt, &req)
	if err != nil {
		fatalf(err.Error())
	}

	workloadID, err := c.CreateWorkload(req)
	if err != nil {
		return errors.Wrap(err, "Error creating workload")
	}

	fmt.Printf("Created new workload: %s\n", workloadID)

	return nil
}

type workloadDeleteCommand struct {
	Flag     flag.FlagSet
	workload string
}

func (cmd *workloadDeleteCommand) usage(...string) {
	fmt.Fprintf(os.Stderr, `usage: ciao-cli [options] workload delete [flags]

Deletes a given workload

The delete flags are:

`)
	cmd.Flag.PrintDefaults()
	os.Exit(2)
}

func (cmd *workloadDeleteCommand) parseArgs(args []string) []string {
	cmd.Flag.StringVar(&cmd.workload, "workload", "", "Workload UUID")
	cmd.Flag.Usage = func() { cmd.usage() }
	cmd.Flag.Parse(args)
	return cmd.Flag.Args()
}

func (cmd *workloadDeleteCommand) run(args []string) error {
	if cmd.workload == "" {
		cmd.usage()
	}

	err := c.DeleteWorkload(cmd.workload)
	if err != nil {
		return errors.Wrap(err, "Error deleting workload")
	}

	return nil
}

type workloadShowCommand struct {
	Flag     flag.FlagSet
	template string
	workload string
}

func (cmd *workloadShowCommand) usage(...string) {
	fmt.Fprintf(os.Stderr, `usage: ciao-cli [options] workload show

Show workload details

`)
	cmd.Flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n%s",
		tfortools.GenerateUsageDecorated("f", types.Workload{}, nil))
	os.Exit(2)
}

func (cmd *workloadShowCommand) parseArgs(args []string) []string {
	cmd.Flag.StringVar(&cmd.workload, "workload", "", "Workload UUID")
	cmd.Flag.StringVar(&cmd.template, "f", "", "Template used to format output")
	cmd.Flag.Usage = func() { cmd.usage() }
	cmd.Flag.Parse(args)
	return cmd.Flag.Args()
}

func (cmd *workloadShowCommand) run(args []string) error {
	var wl types.Workload

	if cmd.workload == "" {
		cmd.usage()
	}

	wl, err := c.GetWorkload(cmd.workload)
	if err != nil {
		return errors.Wrap(err, "Error getting workload")
	}

	if cmd.template != "" {
		return tfortools.OutputToTemplate(os.Stdout, "workload-show", cmd.template, &wl, nil)
	}

	outputWorkload(wl)
	return nil
}

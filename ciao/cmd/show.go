// Copyright © 2017 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/ciao-project/ciao/ciao-sdk"
	"github.com/spf13/cobra"
)

var tenantID string
var computenode string
var detailed bool
var resultlimit int
var marker string
var offset int
var tenant string
var workload string

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show information about various ciao objects",
	Long: `Show outputs a list and/or details for available commands`,
}

var eventShowCmd = &cobra.Command{
	Use:   "event",
	Long: `When called with no args, it will print all events.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(tenantID) != 0 {
			args = append(args, tenantID)
		}
		sdk.Show(cmd, args)
	},
}

var externalipShowCmd = &cobra.Command{
	Use:   "externalip",
	Long: `When called with no args, it will print all externalips.`,
}

var imageShowCmd = &cobra.Command{
	Use:   "image <UUID>",
	Long: `When called with no args, it will print all images.`,
}

var instanceShowCmd = &cobra.Command{
	Use:   "instance <UUID>",
	Long: `When called with no args, it will print all instances.`,
}

var nodeShowCmd = &cobra.Command{
	Use:   "node",
	Long: `When called with no args, it will print all nodes.`,
}

var poolShowCmd = &cobra.Command{
	Use:   "pool [NAME]",
	Long: `When called with no args, it will print all pools.`,
}

var quotasShowCmd = &cobra.Command{
	Use:   "quotas",
	Long: `When called with no args, it will print all quotass.`,
}

var tenantShowCmd = &cobra.Command{
	Use:   "tenant",
	Long: `When called with no args, it will print all tenants.`,
}

var traceShowCmd = &cobra.Command{
	Use:   "trace",
	Long: `When called with no args, it will print all traces.`,
}

var volumeShowCmd = &cobra.Command{
	Use:   "volume",
	Long: `When called with no args, it will print all volumes.`,
}

var workloadShowCmd = &cobra.Command{
	Use:   "workload",
	Long: `When called with no args, it will print all workloads.`,
}
	
var showcmds = []*cobra.Command{eventShowCmd, externalipShowCmd, imageShowCmd, instanceShowCmd, nodeShowCmd, poolShowCmd, quotasShowCmd, tenantShowCmd, traceShowCmd, volumeShowCmd, workloadShowCmd}

func init() {
	for _, cmd := range showcmds {
		// Use the Show API to handle the commands given
		cmd.Run = sdk.Show
		showCmd.AddCommand(cmd)
	}
	RootCmd.AddCommand(showCmd)

	showCmd.PersistentFlags().StringVarP(&sdk.Template, "template", "t", "", "Template used to format output")

	eventShowCmd.Flags().StringVar(&tenantID, "tenant-id", "", "Tenant ID to list events for")

	instanceShowCmd.Flags().StringVar(&computenode, "computenode", "", "Compute node to list instances from (defalut to all  nodes when empty)")
	instanceShowCmd.Flags().BoolVar(&detailed, "verbose", false, "Print detailed information about each instance")
	instanceShowCmd.Flags().IntVar(&resultlimit, "limit", 1, "Limit listing to <limit> results")
	instanceShowCmd.Flags().StringVar(&marker, "marker", "", "Show instance list starting from the next instance after marker")
	instanceShowCmd.Flags().IntVar(&offset, "offset", 0, "Show instance list starting from instance <offset>")
	instanceShowCmd.Flags().StringVar(&tenant, "tenant", "", "Specify to list instances from a tenant other than -tenant-id")
	instanceShowCmd.Flags().StringVar(&workload, "workload", "", "Workload UUID")
}
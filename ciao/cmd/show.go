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
	"strings"

	"github.com/ciao-project/ciao/ciao-sdk"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show information about various ciao objects",
	Long:  `Show outputs a list and/or details for available commands`,
}

var eventShowCmd = &cobra.Command{
	Use:  "event [TENANT ID]",
	Long: `When called with no args, it will print all events.`,
	Run: func(cmd *cobra.Command, args []string) {
		CommandFlags.CommandName = strings.Fields(cmd.Use)[0]
		CommandFlags.Args = args
		sdk.Show(&C, *CommandFlags)
	},
}

var externalipShowCmd = &cobra.Command{
	Use:  "externalip",
	Long: `When called with no args, it will print all externalips.`,
	Run: func(cmd *cobra.Command, args []string) {
		CommandFlags.CommandName = strings.Fields(cmd.Use)[0]
		CommandFlags.Args = args
		sdk.Show(&C, *CommandFlags)
	},
}

var imageShowCmd = &cobra.Command{
	Use:  "image <UUID>",
	Long: `When called with no args, it will print all images.`,
	Run: func(cmd *cobra.Command, args []string) {
		CommandFlags.CommandName = strings.Fields(cmd.Use)[0]
		CommandFlags.Args = args
		sdk.Show(&C, *CommandFlags)
	},
}

var instanceShowCmd = &cobra.Command{
	Use:  "instance <UUID>",
	Long: `When called with no args, it will print all instances.`,
	Run: func(cmd *cobra.Command, args []string) {
		CommandFlags.CommandName = strings.Fields(cmd.Use)[0]
		CommandFlags.Args = args
		sdk.Show(&C, *CommandFlags)
	},
}

var nodeShowCmd = &cobra.Command{
	Use:  "node [NODE-ID]",
	Long: `Show information about a node.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		CommandFlags.CommandName = strings.Fields(cmd.Use)[0]
		CommandFlags.Args = args
		sdk.Show(&C, *CommandFlags)
	},
}

var poolShowCmd = &cobra.Command{
	Use:  "pool [NAME]",
	Long: `Show ciao external IP pool details.`,
	Run: func(cmd *cobra.Command, args []string) {
		CommandFlags.CommandName = strings.Fields(cmd.Use)[0]
		CommandFlags.Args = args
		sdk.Show(&C, *CommandFlags)
	},
}

var quotasShowCmd = &cobra.Command{
	Use:  "quotas [TENANT STRING]",
	Long: `When called with no args, it will print all quotas for current tenant.`,
	Run: func(cmd *cobra.Command, args []string) {
		CommandFlags.CommandName = strings.Fields(cmd.Use)[0]
		CommandFlags.Args = args
		sdk.Show(&C, *CommandFlags)
	},
}

var tenantShowCmd = &cobra.Command{
	Use:  "tenant [NAME]",
	Long: `When called with no args, it will print all tenants.`,
	Run: func(cmd *cobra.Command, args []string) {
		CommandFlags.CommandName = strings.Fields(cmd.Use)[0]
		CommandFlags.Args = args
		sdk.Show(&C, *CommandFlags)
	},
}

var traceShowCmd = &cobra.Command{
	Use:  "trace [LABEL]",
	Long: `When called with no args, it will print all traces.`,
	Run: func(cmd *cobra.Command, args []string) {
		CommandFlags.CommandName = strings.Fields(cmd.Use)[0]
		CommandFlags.Args = args
		sdk.Show(&C, *CommandFlags)
	},
}

var volumeShowCmd = &cobra.Command{
	Use:  "volume [UUID]",
	Long: `When called with no args, it will print all volumes.`,
	Run: func(cmd *cobra.Command, args []string) {
		CommandFlags.CommandName = strings.Fields(cmd.Use)[0]
		CommandFlags.Args = args
		sdk.Show(&C, *CommandFlags)
	},
}

var workloadShowCmd = &cobra.Command{
	Use:  "workload [UUID]",
	Long: `When called with no args, it will print all workloads.`,
	Run: func(cmd *cobra.Command, args []string) {
		CommandFlags.CommandName = strings.Fields(cmd.Use)[0]
		CommandFlags.Args = args
		sdk.Show(&C, *CommandFlags)
	},
}

var showcmds = []*cobra.Command{eventShowCmd, externalipShowCmd, imageShowCmd, instanceShowCmd, nodeShowCmd, poolShowCmd, quotasShowCmd, tenantShowCmd, traceShowCmd, volumeShowCmd, workloadShowCmd}

func init() {
	for _, cmd := range showcmds {
		showCmd.AddCommand(cmd)
	}

	RootCmd.AddCommand(showCmd)

	showCmd.PersistentFlags().StringVarP(&C.Template, "template", "t", "", "Template used to format output")

	eventShowCmd.Flags().BoolVar(&CommandFlags.All, "all", false, "List events for all tenants in a cluster")
	eventShowCmd.Flags().StringVar(&CommandFlags.Tenant, "tenant", "", "Tenant to list events from")

	instanceShowCmd.Flags().StringVar(&CommandFlags.Computenode, "computenode", "", "Compute node to list instances from (defalut to all  nodes when empty)")
	instanceShowCmd.Flags().BoolVar(&CommandFlags.Detail, "verbose", false, "Print detailed information about each instance")
	instanceShowCmd.Flags().IntVar(&CommandFlags.Limit, "limit", 1, "Limit listing to <limit> results")
	instanceShowCmd.Flags().StringVar(&CommandFlags.Marker, "marker", "", "Show instance list starting from the next instance after marker")
	instanceShowCmd.Flags().IntVar(&CommandFlags.Offset, "offset", 0, "Show instance list starting from instance <offset>")
	instanceShowCmd.Flags().StringVar(&CommandFlags.Tenant, "tenant", "", "Specify to list instances from a tenant other than -tenant-id")
	instanceShowCmd.Flags().StringVar(&CommandFlags.Workload, "workload", "", "Workload UUID")
}

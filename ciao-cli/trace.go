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

var traceCommand = &command{
	SubCommands: map[string]subCommand{
		"list": new(traceListCommand),
		"show": new(traceShowCommand),
	},
}

type traceListCommand struct {
	Flag     flag.FlagSet
	template string
}

func (cmd *traceListCommand) usage(...string) {
	fmt.Fprintf(os.Stderr, `usage: ciao-cli [options] trace list

List all trace label
`)
	cmd.Flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n%s",
		tfortools.GenerateUsageDecorated("f", types.CiaoTracesSummary{}.Summaries, nil))
	os.Exit(2)
}

func (cmd *traceListCommand) parseArgs(args []string) []string {
	cmd.Flag.StringVar(&cmd.template, "f", "", "Template used to format output")
	cmd.Flag.Usage = func() { cmd.usage() }
	cmd.Flag.Parse(args)
	return cmd.Flag.Args()
}

func (cmd *traceListCommand) run(args []string) error {
	traces, err := client.ListTraceLabels()
	if err != nil {
		return errors.Wrap(err, "Error listing trace labels")
	}

	if cmd.template != "" {
		return tfortools.OutputToTemplate(os.Stdout, "trace-list", cmd.template,
			&traces.Summaries, nil)
	}

	fmt.Printf("%d trace label(s) available\n", len(traces.Summaries))
	for i, summary := range traces.Summaries {
		fmt.Printf("\tLabel #%d: %s (%d instances running)\n", i+1, summary.Label, summary.Instances)
	}

	return nil
}

type traceShowCommand struct {
	Flag     flag.FlagSet
	label    string
	template string
}

func (cmd *traceShowCommand) usage(...string) {
	fmt.Fprintf(os.Stderr, `usage: ciao-cli [options] trace show [flags]

Dump all trace data for a given label

The show flags are:

`)
	cmd.Flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n%s",
		tfortools.GenerateUsageDecorated("f", types.CiaoTraceData{}.Summary, nil))
	os.Exit(2)
}

func (cmd *traceShowCommand) parseArgs(args []string) []string {
	cmd.Flag.StringVar(&cmd.label, "label", "", "Label name")
	cmd.Flag.StringVar(&cmd.template, "f", "", "Template used to format output")
	cmd.Flag.Usage = func() { cmd.usage() }
	cmd.Flag.Parse(args)
	return cmd.Flag.Args()
}

func (cmd *traceShowCommand) run(args []string) error {
	if cmd.label == "" {
		return errors.New("Missing required -label parameter")
	}

	traceData, err := client.GetTraceData(cmd.label)
	if err != nil {
		return errors.Wrap(err, "Error getting trace data")
	}

	if cmd.template != "" {
		return tfortools.OutputToTemplate(os.Stdout, "trace-show", cmd.template,
			&traceData.Summary, nil)
	}

	fmt.Printf("Trace data for [%s]:\n", cmd.label)
	fmt.Printf("\tNumber of instances: %d\n", traceData.Summary.NumInstances)
	fmt.Printf("\tTotal time elapsed     : %f seconds\n", traceData.Summary.TotalElapsed)
	fmt.Printf("\tAverage time elapsed   : %f seconds\n", traceData.Summary.AverageElapsed)
	fmt.Printf("\tAverage Controller time: %f seconds\n", traceData.Summary.AverageControllerElapsed)
	fmt.Printf("\tAverage Scheduler time : %f seconds\n", traceData.Summary.AverageSchedulerElapsed)
	fmt.Printf("\tAverage Launcher time  : %f seconds\n", traceData.Summary.AverageLauncherElapsed)
	fmt.Printf("\tController variance    : %f seconds²\n", traceData.Summary.VarianceController)
	fmt.Printf("\tScheduler variance     : %f seconds²\n", traceData.Summary.VarianceScheduler)
	fmt.Printf("\tLauncher variance      : %f seconds²\n", traceData.Summary.VarianceLauncher)

	return nil
}

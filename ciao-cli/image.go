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

var imageCommand = &command{
	SubCommands: map[string]subCommand{
		"add":    new(imageAddCommand),
		"show":   new(imageShowCommand),
		"list":   new(imageListCommand),
		"delete": new(imageDeleteCommand),
	},
}

type imageAddCommand struct {
	Flag       flag.FlagSet
	name       string
	id         string
	file       string
	template   string
	visibility string
}

func (cmd *imageAddCommand) usage(...string) {
	fmt.Fprintf(os.Stderr, `usage: ciao-cli [options] image add [flags]

Creates a new image

The add flags are:

`)
	cmd.Flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n%s", tfortools.GenerateUsageDecorated("f", types.Image{}, nil))
	os.Exit(2)
}

func (cmd *imageAddCommand) parseArgs(args []string) []string {
	cmd.Flag.StringVar(&cmd.name, "name", "", "Image Name")
	cmd.Flag.StringVar(&cmd.id, "id", "", "Image UUID")
	cmd.Flag.StringVar(&cmd.file, "file", "", "Image file to upload")
	cmd.Flag.StringVar(&cmd.template, "f", "", "Template used to format output")
	cmd.Flag.StringVar(&cmd.visibility, "visibility", string(types.Private),
		"Image visibility (internal,public,private)")
	cmd.Flag.Usage = func() { cmd.usage() }
	cmd.Flag.Parse(args)
	return cmd.Flag.Args()
}

func (cmd *imageAddCommand) run(args []string) error {
	if cmd.name == "" {
		return errors.New("Missing required -name parameter")
	}

	if cmd.file == "" {
		return errors.New("Missing required -file parameter")
	}

	f, err := os.Open(cmd.file)
	if err != nil {
		fatalf("Could not open %s [%s]\n", cmd.file, err)
	}
	defer func() { _ = f.Close() }()

	imageVisibility := types.Private
	if cmd.visibility != "" {
		imageVisibility = types.Visibility(cmd.visibility)
		switch imageVisibility {
		case types.Public, types.Private, types.Internal:
		default:
			fatalf("Invalid image visibility [%v]", imageVisibility)
		}
	}

	id, err := c.CreateImage(cmd.name, imageVisibility, cmd.id, f)
	if err != nil {
		return errors.Wrap(err, "Error creating image")
	}

	image, err := c.GetImage(id)
	if err != nil {
		return errors.Wrap(err, "Error getting image")
	}

	if cmd.template != "" {
		return tfortools.OutputToTemplate(os.Stdout, "image-add", cmd.template, image, nil)
	}

	fmt.Printf("Created image:\n")
	dumpImage(&image)
	return nil
}

type imageShowCommand struct {
	Flag     flag.FlagSet
	image    string
	template string
}

func (cmd *imageShowCommand) usage(...string) {
	fmt.Fprintf(os.Stderr, `usage: ciao-cli [options] image show

Show images
`)
	cmd.Flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n%s", tfortools.GenerateUsageDecorated("f", types.Image{}, nil))
	os.Exit(2)
}

func (cmd *imageShowCommand) parseArgs(args []string) []string {
	cmd.Flag.StringVar(&cmd.template, "f", "", "Template used to format output")
	cmd.Flag.StringVar(&cmd.image, "image", "", "Image UUID")
	cmd.Flag.Usage = func() { cmd.usage() }
	cmd.Flag.Parse(args)
	return cmd.Flag.Args()
}

func (cmd *imageShowCommand) run(args []string) error {
	if cmd.image == "" {
		return errors.New("Missing required -image parameter")
	}

	i, err := c.GetImage(cmd.image)
	if err != nil {
		return errors.Wrap(err, "Error getting image")
	}

	if cmd.template != "" {
		return tfortools.OutputToTemplate(os.Stdout, "image-show", cmd.template, i, nil)
	}

	dumpImage(&i)

	return nil
}

type imageListCommand struct {
	Flag     flag.FlagSet
	template string
}

func (cmd *imageListCommand) usage(...string) {
	fmt.Fprintf(os.Stderr, `usage: ciao-cli [options] image list

List images
`)
	cmd.Flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
The template passed to the -f option operates on a 

%s

As images are retrieved in pages, the template may be applied multiple
times.  You can not therefore rely on the length of the slice passed
to the template to determine the total number of images.
`, tfortools.GenerateUsageUndecorated([]types.Image{}))
	fmt.Fprintln(os.Stderr, tfortools.TemplateFunctionHelp(nil))
	os.Exit(2)
}

func (cmd *imageListCommand) parseArgs(args []string) []string {
	cmd.Flag.StringVar(&cmd.template, "f", "", "Template used to format output")
	cmd.Flag.Usage = func() { cmd.usage() }
	cmd.Flag.Parse(args)
	return cmd.Flag.Args()
}

func (cmd *imageListCommand) run(args []string) error {
	err := c.ListImages()
	if err != nil {
		return errors.Wrap(err, "Error listing images")
	}

	return err
}

type imageDeleteCommand struct {
	Flag  flag.FlagSet
	image string
}

func (cmd *imageDeleteCommand) usage(...string) {
	fmt.Fprintf(os.Stderr, `usage: ciao-cli [options] image delete [flags]

Deletes an image

The delete flags are:

`)
	cmd.Flag.PrintDefaults()
	os.Exit(2)
}

func (cmd *imageDeleteCommand) parseArgs(args []string) []string {
	cmd.Flag.StringVar(&cmd.image, "image", "", "Image UUID")
	cmd.Flag.Usage = func() { cmd.usage() }
	cmd.Flag.Parse(args)
	return cmd.Flag.Args()
}

func (cmd *imageDeleteCommand) run(args []string) error {
	err := c.DeleteImage(cmd.image)
	if err != nil {
		return errors.Wrap(err, "Error deleting image")
	}

	fmt.Printf("Deleted image %s\n", cmd.image)

	return nil
}

func dumpImage(i *types.Image) {
	fmt.Printf("\tName\t\t[%s]\n", i.Name)
	fmt.Printf("\tSize\t\t[%d bytes]\n", i.Size)
	fmt.Printf("\tID\t\t[%s]\n", i.ID)
	fmt.Printf("\tState\t\t[%s]\n", i.State)
	fmt.Printf("\tVisibility\t[%s]\n", i.Visibility)
	fmt.Printf("\tCreateTime\t[%s]\n", i.CreateTime)
}

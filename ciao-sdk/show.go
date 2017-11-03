package sdk

import (
	"github.com/ciao-project/ciao/client"

	"github.com/pkg/errors"
)

func Show(c *client.Client, data CommandOpts) {
	var ret error

	switch data.CommandName {
	case "instance":
		if len(data.Args) == 0 {
			ret = ListInstances(c, data)
		} else {
			ret = ShowInstance(c, data)
		}
	case "image":
		if len(data.Args) == 0 {
			c.ListImages()
		} else {
			c.ShowImage(data.Args[0])
		}
	}
	if ret != nil {
		errors.Wrapf(ret, "Error running %s\n", data.CommandName)
	}
}

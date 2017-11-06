package sdk

import (
	"github.com/ciao-project/ciao/client"

	"github.com/pkg/errors"
)

func Show(c *client.Client, data CommandOpts) {
	var ret error

	switch data.CommandName {
	case "event":
		ret = ListEvents(c, data)
	case "externalip":
		ret = ListExternalIP(c, data)
	case "instance":
		if len(data.Args) == 0 {
			ret = ListInstances(c, data)
		} else {
			ret = ShowInstance(c, data)
		}
	case "image":
		if len(data.Args) == 0 {
			ret = c.ListImages()
		} else {
			ret = c.ShowImage(data.Args[0])
		}
	case "workload":
		if len(data.Args) == 0 {
			ret = ListWorkload(c, data)
		} else {
			ret = ShowWorkload(c, data)
		}
	}
	if ret != nil {
		errors.Wrapf(ret, "Error running %s\n", data.CommandName)
	}
}

package sdk

import (
	"bytes"

	"github.com/ciao-project/ciao/ciao-controller/types"
	"github.com/ciao-project/ciao/client"

	"github.com/pkg/errors"
)

/* Intended calling convention by user would be:
 * instances := sdk.Show(ClientHandle, "instance", nil)
 * for example, to return a string of all instances using
 * the default tfortools template {{table .}}. This can
 * then be printed or parsed as needed. */
func Show(c *client.Client, objName string, data CommandOpts) (bytes.Buffer, error) {
	var err error
	var result bytes.Buffer

	switch objName {
	case "event":
		events, err := ListEvents(c, data)
		if err == nil {
			c.PrettyPrint(&result, "list-events", events)
		}
	case "externalip":
		IPs, err := ListExternalIP(c, data)
		if err == nil {
			c.PrettyPrint(&result, "list-externalip", IPs)
		}
	case "instance":
		if len(data.Args) == 0 {
			instances, err := ListInstances(c, data)
			if err == nil {
				c.PrettyPrint(&result, "list-instance", instances)
			}
		} else {
			instance, err := ShowInstance(c, data)
			if err == nil {
				c.PrettyPrint(&result, "list-instance", instance)
			}
		}
	case "image":
		if len(data.Args) == 0 {
			images, err := GetImageList(c, data)
			if err == nil {
				c.PrettyPrint(&result, "list-image", images)
			}
		} else {
			image, err := GetImage(c, data)
			images := []types.Image{image}
			if err == nil {
				c.PrettyPrint(&result, "list-image", images)
			}
		}
	case "workload":
		if len(data.Args) == 0 {
			workloads, err := GetWorkloadList(c, data)
			if err == nil {
				c.PrettyPrint(&result, "list-workload", workloads)
			}
		} else {
			workload, err := GetWorkload(c, data)
			if err == nil {
				wl := []Workload{workload}
				c.PrettyPrint(&result, "list-workload", wl)
			}
		}
	}
	if err != nil {
		return result, errors.Wrapf(err, "Error running %s\n", objName)
	}

	return result, nil
}

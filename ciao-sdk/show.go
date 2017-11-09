package sdk

import (
	"bytes"

	"github.com/ciao-project/ciao/client"
	"github.com/ciao-project/ciao/ciao-controller/types"

	"github.com/intel/tfortools"

	"github.com/pkg/errors"
)

/* Intended calling convention by user would be:
 * instances := sdk.Show(ClientHandle, "instance", nil)
 * for example, to return a string of all instances using
 * the default tfortools template {{table .}}. This can
 * then be printed or parsed as needed. */
func Show(c *client.Client, objName string, data CommandOpts) (string, error) {
	var err error
	var result bytes.Buffer

	switch objName {
	case "event":
		events, err := ListEvents(c, data)
		if err == nil {
			tfortools.OutputToTemplate(&result, "event-list", "{{table .}}", events, nil)
		}
	case "externalip":
		IPs, err := ListExternalIP(c, data)
		if err == nil {
			tfortools.OutputToTemplate(&result, "externalip-list", "{{table .}}", IPs, nil)
		}
	case "instance":
		if len(data.Args) == 0 {
			instances, err := ListInstances(c, data)
			if err == nil {
				tfortools.OutputToTemplate(&result, "instance-list", "{{table .}}", instances, nil)
			}
		} else {
			instance, err := ShowInstance(c, data)
			if err == nil {
				tfortools.OutputToTemplate(&result, "instance-show", "{{table .}}", instance, nil)
			}
		}
	case "image":
		if len(data.Args) == 0 {
			images, err := GetImageList(c, data)
			if err == nil {
				tfortools.OutputToTemplate(&result, "workload-show", "{{table .}}", images, nil)
			}
		} else {
			image, err := GetImage(c, data)
			images := []types.Image{image}
			if err == nil {
				tfortools.OutputToTemplate(&result, "workload-show", "{{table .}}", images, nil)
			}
		}

	case "workload":
		if len(data.Args) == 0 {
			workloads, err := GetWorkloadList(c, data)
			if err == nil {
				tfortools.OutputToTemplate(&result, "workload-list", "{{table .}}", workloads, nil)
			}
		} else {
			workload, err := GetWorkload(c, data)
			if err == nil {
				wl := []Workload{workload}
				tfortools.OutputToTemplate(&result, "workload-show", "{{table .}}", wl, nil)
			}
		}
	}
	if err != nil {
		return "", errors.Wrapf(err, "Error running %s\n", objName)
	}

	return result.String(), nil
}

package sdk

import (
	"github.com/ciao-project/ciao/client"

	"github.com/pkg/errors"

)
var c client.Client

func Show(object string, args []string) {
	var ret error

	switch object {

	case "image":
		c.ListImages()
	}
	if ret != nil {
		errors.Wrapf(ret, "Error running %s\n", object)
	}
}

package sdk

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/ciao-project/ciao/client"
	"github.com/ciao-project/ciao/ciao-controller/api"
	"github.com/intel/tfortools"
	"github.com/pkg/errors"
	"github.com/golang/glog"

)

func fatalf(format string, args ...interface{}) {
	glog.FatalDepth(1, fmt.Sprintf("ciao-cli FATAL: "+format, args...))
}


type byCreated []api.ServerDetails

func (ss byCreated) Len() int           { return len(ss) }
func (ss byCreated) Swap(i, j int)      { ss[i], ss[j] = ss[j], ss[i] }
func (ss byCreated) Less(i, j int) bool { return ss[i].Created.Before(ss[j].Created) }

func dumpInstance(server *api.ServerDetails) {
	fmt.Printf("\tUUID: %s\n", server.ID)
	fmt.Printf("\tStatus: %s\n", server.Status)
	fmt.Printf("\tPrivate IP: %s\n", server.PrivateAddresses[0].Addr)
	fmt.Printf("\tMAC Address: %s\n", server.PrivateAddresses[0].MacAddr)
	fmt.Printf("\tCN UUID: %s\n", server.NodeID)
	fmt.Printf("\tTenant UUID: %s\n", server.TenantID)
	if server.SSHIP != "" {
		fmt.Printf("\tSSH IP: %s\n", server.SSHIP)
		fmt.Printf("\tSSH Port: %d\n", server.SSHPort)
	}

	for _, vol := range server.Volumes {
		fmt.Printf("\tVolume: %s\n", vol)
	}
}

func listNodeInstances(c *client.Client, node string) error {
	if node == "" {
		fatalf("Missing required -cn parameter")
	}

	servers, err := c.ListInstancesByNode(node)
	if err != nil {
		return errors.Wrap(err, "Error getting instances for node")
	}

	for i, server := range servers.Servers {
		fmt.Printf("Instance #%d\n", i+1)
		fmt.Printf("\tUUID: %s\n", server.ID)
		fmt.Printf("\tStatus: %s\n", server.Status)
		fmt.Printf("\tTenant UUID: %s\n", server.TenantID)
		fmt.Printf("\tIPv4: %s\n", server.IPv4)
		fmt.Printf("\tCPUs used: %d\n", server.VCPUUsage)
		fmt.Printf("\tMemory used: %d MB\n", server.MemUsage)
		fmt.Printf("\tDisk used: %d MB\n", server.DiskUsage)
	}

	return nil
}

func ListInstances(c *client.Client, flags CommandOpts) error {
	if flags.Tenant == "" {
		flags.Tenant = c.TenantID
	}

	if flags.Computenode != "" {
		return listNodeInstances(c, flags.Computenode)
	}

	servers, err := c.ListInstancesByWorkload(flags.Tenant, flags.Workload)
	if err != nil {
		return errors.Wrap(err, "Error listing instances")
	}

	sortedServers := []api.ServerDetails{}
	for _, v := range servers.Servers {
		sortedServers = append(sortedServers, v)
	}
	sort.Sort(byCreated(sortedServers))

	if c.Template != "" {
		return tfortools.OutputToTemplate(os.Stdout, "instance-list", c.Template,
			&sortedServers, nil)
	}

	w := new(tabwriter.Writer)
	if !flags.Detail {
		w.Init(os.Stdout, 0, 1, 1, ' ', 0)
		fmt.Fprintln(w, "#\tUUID\tStatus\tPrivate IP\tSSH IP\tSSH PORT")
	}

	for i, server := range sortedServers {
		if !flags.Detail {
			fmt.Fprintf(w, "%d", i+1)
			fmt.Fprintf(w, "\t%s", server.ID)
			fmt.Fprintf(w, "\t%s", server.Status)
			fmt.Fprintf(w, "\t%s", server.PrivateAddresses[0].Addr)
			if server.SSHIP != "" {
				fmt.Fprintf(w, "\t%s", server.SSHIP)
				fmt.Fprintf(w, "\t%d\n", server.SSHPort)
			} else {
				fmt.Fprintf(w, "\tN/A")
				fmt.Fprintf(w, "\tN/A\n")
			}
			w.Flush()
		} else {
			fmt.Printf("Instance #%d\n", i+1)
			dumpInstance(&server)
		}
	}
	return nil
}
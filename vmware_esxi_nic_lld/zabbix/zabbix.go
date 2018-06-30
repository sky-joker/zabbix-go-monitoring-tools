package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/urfave/cli"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type ZabbixLLD struct {
	Data []interface{} `json:"data"`
}

func doMain(opt *cli.Context) {
	app := doApp()
	v := opt.Args().Get(0)
	usr := opt.Args().Get(1)
	pwd := opt.Args().Get(2)
	uuid := opt.Args().Get(3)
	datacenter := opt.Args().Get(4)
	if v == "" || usr == "" || pwd == "" || uuid == "" {
		help := []string{"", "--help"}
		app.Run(help)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	u, err := url.Parse(v)
	if err != nil {
		os.Exit(1)
	}

	u.User = url.UserPassword(usr, pwd)
	if err != nil {
		os.Exit(1)
	}

	c, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		os.Exit(1)
	}

	f := find.NewFinder(c.Client, true)
	dc, _ := f.Datacenter(ctx, "/"+datacenter)

	s := object.NewSearchIndex(c.Client)
	r, _ := s.FindByUuid(ctx, dc, uuid, false, nil)

	var host []mo.HostSystem
	var refs = []types.ManagedObjectReference{r.Reference()}
	pc := property.DefaultCollector(c.Client)
	pc.Retrieve(ctx, refs, []string{"name", "config"}, &host)

	zbx := new(ZabbixLLD)
	var Nics struct {
		Nic string `json:"{#NIC}"`
	}
	for _, v := range host {
		for _, p := range v.Config.Network.Pnic {
			Nics.Nic = p.Device
			zbx.Data = append(zbx.Data, Nics)
		}
	}

	j, err := json.Marshal(zbx)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(string(j))
}

func doApp() *cli.App {
	app := cli.NewApp()
	app.Name = "vmware_esxi_nic_lld"
	app.Usage = "Get physical nic name of esxi"
	app.UsageText = "vmware_esxi_nic_lld URL UserName Password Uuid Datacenter"
	app.Commands = commands
	app.HideVersion = true
	return app
}

func Do() {
	app := doApp()
	app.Action = doMain
	app.Run(os.Args)
}

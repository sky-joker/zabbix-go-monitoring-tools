package zabbix

import "github.com/urfave/cli"

var commands = []cli.Command{
	URL,
	UserName,
	Password,
	Uuid,
	Datacenter,
}

var URL = cli.Command{
	Name:  "URL",
	Usage: "URL for vCenter sdk(ex:https://vcenter/sdk)",
}

var UserName = cli.Command{
	Name:  "UserName",
	Usage: "vCenter user",
}

var Password = cli.Command{
	Name:  "Password",
	Usage: "vCenter user password",
}

var Uuid = cli.Command{
	Name:  "Uuid",
	Usage: "ESXi uuid",
}

var Datacenter = cli.Command{
	Name:  "Datacenter",
	Usage: "Datacenter name(optional)",
}

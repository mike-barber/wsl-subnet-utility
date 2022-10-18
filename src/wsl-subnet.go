package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Microsoft/hcsshim/hcn"
	log "github.com/sirupsen/logrus"
	"net/netip"
)

func ExistsWsl() bool {
	networks, err := hcn.ListNetworks()
	if err != nil {
		log.Fatal(err)
	}
	for _, n := range networks {
		if n.Name == "WSL" {
			return true
		}
	}
	return false
}

func DeleteWsl() {
	log.Info("Removing old WSL network...")

	orig, err := hcn.GetNetworkByName("WSL")
	if err != nil {
		log.Fatal(err)
	}

	errDel := orig.Delete()
	if errDel != nil {
		log.Fatal(errDel)
	}
}

// This appears to work correctly.
func CreateWsl(prefix netip.Prefix) {
	log.Info("Creating new WSL network...")
	ipHost := prefix.Addr().Next()

	newNet := hcn.HostComputeNetwork{
		Name: "WSL",
		Type: "ICS",
		Id:   "B95D0C5E-57D4-412B-B571-18A81A16E005",
		Dns: hcn.Dns{
			ServerList: []string{
				ipHost.String(),
			},
		},
		Ipams: []hcn.Ipam{
			{
				Subnets: []hcn.Subnet{
					{
						IpAddressPrefix: prefix.String(),
						Routes: []hcn.Route{
							{
								NextHop:           ipHost.String(),
								DestinationPrefix: "0.0.0.0/0",
							},
						},
					},
				},
			},
		},
		Flags:         41,
		SchemaVersion: hcn.Version{Major: 2, Minor: 0},
	}

	resultNet, err := newNet.Create()
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{
		"ip":   ipHost.String(),
		"net":  prefix.String(),
		"guid": resultNet.Id,
	}).Info("New WSL network created successfully")
}

func PrintNetwork(v interface{}) {
	bytes, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(bytes))
}

func ListNetworks() {
	networks, err := hcn.ListNetworks()
	if err != nil {
		log.Fatal(err)
	}
	for _, n := range networks {
		PrintNetwork(n)
	}
}

func main() {
	var fDelete = flag.Bool("allow-delete", true, "allow deletion if existing WSL network")
	var fCreate = flag.Bool("allow-create", true, "allow creation of network after deletion")
	var fList = flag.Bool("list", false, "list all wsl networks before and after changes")
	var fIpSubnet = flag.String("ipnet", "192.168.100.0/24", "WSL host subnet; must contain host address above")
	flag.Parse()

	prefix, err := netip.ParsePrefix(*fIpSubnet)
	if err != nil || !prefix.Addr().Is4() {
		log.Fatal("You didn't specify a valid IPv4 prefix")
	}

	if *fList {
		log.Info("Networks before action:")
		ListNetworks()
	}

	if ExistsWsl() {
		if *fDelete {
			DeleteWsl()
		} else {
			log.Fatal("WSL network already exists and allow-delete is not enabled")
		}
	} else {
		log.Info("No pre-existing WSL network")
	}

	if *fCreate {
		CreateWsl(prefix)
	}

	if *fList {
		log.Info("Networks after action:")
		ListNetworks()
	}
}

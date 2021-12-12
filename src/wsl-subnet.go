package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Microsoft/hcsshim/hcn"
	log "github.com/sirupsen/logrus"
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
func CreateWsl(ipHost string, ipSubnet string) {
	log.Info("Creating new WSL network...")

	newNet := hcn.HostComputeNetwork{
		Name: "WSL",
		Type: "ICS",
		Id:   "B95D0C5E-57D4-412B-B571-18A81A16E005",
		MacPool: hcn.MacPool{
			Ranges: []hcn.MacRange{
				{
					StartMacAddress: "00-15-5D-EC-B0-00",
					EndMacAddress:   "00-15-5D-EC-BF-FF",
				},
			},
		},
		Dns: hcn.Dns{
			ServerList: []string{
				ipHost,
			},
		},
		Ipams: []hcn.Ipam{
			{
				Subnets: []hcn.Subnet{
					{
						IpAddressPrefix: ipSubnet,
						Routes: []hcn.Route{
							{
								NextHop:           ipHost,
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
		"ip":   ipHost,
		"net":  ipSubnet,
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
	var fIpHost = flag.String("ip", "192.168.100.1", "WSL host machine IP address")
	var fIpSubnet = flag.String("ipnet", "192.168.100.0/24", "WSL host subnet; must contain host address above")
	flag.Parse()

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
		CreateWsl(*fIpHost, *fIpSubnet)
	}

	if *fList {
		log.Info("Networks after action:")
		ListNetworks()
	}
}

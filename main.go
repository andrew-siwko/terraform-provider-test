package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	// print_macs()
	// newPrintMacs()
	// fmt.Print("Starting terraform-provider-test\n")
	err := providerserver.Serve(context.Background(), New, providerserver.ServeOpts{
		Address: "registry.terraform.io/andrew/property-mirror",
	})

	if err != nil {
		// fmt.Print("Ending terraform-provider-test\n")
		// maybe we coult test the web service here.
		log.Fatal(err.Error())
	}
	// fmt.Print("Ending terraform-provider-test\n")
}

func NewPrintMacs() {
	fmt.Println("Detecting primary outbound MAC addresses...")
	fmt.Println("------------------------------------------")

	macList, err := getLocalMACs()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Iterate and print each discovered MAC address
	for i, mac := range macList {
		fmt.Printf("[%d] MAC Address: %s\n", i+1, mac)
	}
}

func Print_macs() {
	// Retrieve all network interfaces on the system
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Error fetching interfaces: %v\n", err)
		return
	}

	fmt.Println("Available MAC Addresses:")
	fmt.Println("------------------------")

	for _, iface := range interfaces {
		// Skip interfaces that are down or have no hardware address (like loopback)
		if len(iface.HardwareAddr) == 0 {
			continue
		}

		// Optional: Filter out down interfaces if you only want active connections
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		fmt.Printf("Interface: %-15s MAC: %s\n", iface.Name, iface.HardwareAddr.String())
	}
}

func getLocalMACs() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch interfaces: %w", err)
	}
	var macs []string

	for _, iface := range interfaces {
		// Skip interfaces with no hardware address (like loopback)
		if len(iface.HardwareAddr) == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if _, ok := addr.(*net.IPNet); ok {
				macs = append(macs, iface.HardwareAddr.String())
				break // Move to the next interface if this one matches
			}
		}
	}

	return macs, nil
}

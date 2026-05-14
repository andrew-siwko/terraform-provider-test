package main

// this did not work
// gowsdl -p vboxapi http://daddy.siwko.org:18083/?wsdl > /vboxapi/vbox_bindings.go

import (
	"github.com/andrew/terraform-provider-test/vboxapi" // The generated code
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hooklift/gowsdl/soap"
)

func (c *VBoxClient) GetVMNames() ([]string, error) {
	soapClient := soap.NewClient(c.Endpoint)
	service := vboxapi.NewVboxPortType(soapClient)

	resp, err := service.IWebsessionManager_logon(&vboxapi.IWebsessionManager_logon{
		Username: c.Username, // Can be ""
		Password: c.Password, // Can be ""
	})
	if err != nil {
		return nil, err
	}

	vboxHandle := resp.Returnval

	machinesResp, err := service.IVirtualBox_getMachines(&vboxapi.IVirtualBox_getMachines{
		This: vboxHandle,
	})
	if err != nil {
		return nil, err
	}

	var names []string
	for _, machineHandle := range machinesResp.Returnval {
		nameResp, err := service.IMachine_getName(&vboxapi.IMachine_getName{
			This: machineHandle,
		})
		if err == nil {
			names = append(names, nameResp.Returnval)
		}
	}

	service.IWebsessionManager_logoff(&vboxapi.IWebsessionManager_logoff{
		RefIVirtualBox: vboxHandle,
	})

	return names, nil
}

func (c *VBoxClient) GetDetailedVMs() ([]vmModel, error) {
	soapClient := soap.NewClient(c.Endpoint)
	service := vboxapi.NewVboxPortType(soapClient)

	resp, err := service.IWebsessionManager_logon(&vboxapi.IWebsessionManager_logon{
		Username: c.Username, // Can be ""
		Password: c.Password, // Can be ""
	})
	if err != nil {
		return nil, err
	}

	vboxHandle := resp.Returnval

	// Get Machine Handles
	machinesResp, err := service.IVirtualBox_getMachines(&vboxapi.IVirtualBox_getMachines{
		This: vboxHandle,
	})
	if err != nil {
		return nil, err
	}

	var results []vmModel

	for _, handle := range machinesResp.Returnval {
		// Get Name
		n, _ := service.IMachine_getName(&vboxapi.IMachine_getName{This: handle})
		
		// Get Memory (VirtualBox returns MB)
		m, _ := service.IMachine_getMemorySize(&vboxapi.IMachine_getMemorySize{This: handle})
		
		// Get CPU Count
		cp, _ := service.IMachine_getCPUCount(&vboxapi.IMachine_getCPUCount{This: handle})
		
		// Get State
		s, _ := service.IMachine_getState(&vboxapi.IMachine_getState{This: handle})

		var stateString string
		if s != nil && s.Returnval != nil {
			stateString = string(*s.Returnval)
		} else {
			stateString = "Unknown"
		}


		results = append(results, vmModel{
			Name:   types.StringValue(n.Returnval),
			Memory: types.Int64Value(int64(m.Returnval)),
			CPUs:   types.Int64Value(int64(cp.Returnval)),
			State:  types.StringValue(string(stateString)),
		})
	}

	// Always Logoff!
	service.IWebsessionManager_logoff(&vboxapi.IWebsessionManager_logoff{
		RefIVirtualBox: vboxHandle,
	})

	return results, nil
}
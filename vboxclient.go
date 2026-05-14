package main

// this did not work
// gowsdl -p vboxapi http://daddy.siwko.org:18083/?wsdl > internal/vboxapi/vbox_bindings.go

import (
	"github.com/andrew/terraform-provider-test/vboxapi" // The generated code
	"github.com/hooklift/gowsdl/soap"
)

func (c *VBoxClient) GetVMNames() ([]string, error) {
	// 1. Create the SOAP client
	soapClient := soap.NewClient(c.Endpoint)
	service := vboxapi.NewVboxPortType(soapClient)

	// 2. "Logon" to get the IVirtualBox handle
	// Even unauthenticated, we need this handle to talk to the hypervisor
	resp, err := service.IWebsessionManager_logon(&vboxapi.IWebsessionManager_logon{
		Username: c.Username, // Can be ""
		Password: c.Password, // Can be ""
	})
	if err != nil {
		return nil, err
	}
	vboxHandle := resp.Returnval

	// 3. Get the list of Machine handles
	machinesResp, err := service.IVirtualBox_getMachines(&vboxapi.IVirtualBox_getMachines{
		This: vboxHandle,
	})
	if err != nil {
		return nil, err
	}

	// 4. Loop through handles to get names
	var names []string
	for _, machineHandle := range machinesResp.Returnval {
		nameResp, err := service.IMachine_getName(&vboxapi.IMachine_getName{
			This: machineHandle,
		})
		if err == nil {
			names = append(names, nameResp.Returnval)
		}
	}

	// 5. Always logoff to clear the session on your RHEL host
	service.IWebsessionManager_logoff(&vboxapi.IWebsessionManager_logoff{
		This: vboxHandle,
	})

	return names, nil
}

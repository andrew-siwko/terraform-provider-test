package main

// this did not work
// gowsdl -p vboxapi http://daddy.siwko.org:18083/?wsdl > /vboxapi/vbox_bindings.go

import (
	"github.com/andrew/terraform-provider-test/vboxapi" // The generated code
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

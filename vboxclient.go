package main

// this did not work
// gowsdl -p vboxapi http://daddy.siwko.org:18083/?wsdl > /vbox_interface.go

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hooklift/gowsdl/soap"
)

func (c *VBoxClient) GetVMNames() ([]string, error) {
	soapClient := soap.NewClient(c.Endpoint)
	service := NewVboxPortType(soapClient)

	resp, err := service.IWebsessionManager_logon(&IWebsessionManager_logon{
		Username: c.Username, // Can be ""
		Password: c.Password, // Can be ""
	})
	if err != nil {
		return nil, err
	}

	vboxHandle := resp.Returnval

	machinesResp, err := service.IVirtualBox_getMachines(&IVirtualBox_getMachines{
		This: vboxHandle,
	})
	if err != nil {
		return nil, err
	}

	var names []string
	for _, machineHandle := range machinesResp.Returnval {
		nameResp, err := service.IMachine_getName(&IMachine_getName{
			This: machineHandle,
		})
		if err == nil {
			names = append(names, nameResp.Returnval)
		}
	}

	service.IWebsessionManager_logoff(&IWebsessionManager_logoff{
		RefIVirtualBox: vboxHandle,
	})

	return names, nil
}

func (c *VBoxClient) GetDetailedVMs() ([]vmModel, error) {
	soapClient := soap.NewClient(c.Endpoint)
	service := NewVboxPortType(soapClient)

	resp, err := service.IWebsessionManager_logon(&IWebsessionManager_logon{
		Username: c.Username, // Can be ""
		Password: c.Password, // Can be ""
	})
	if err != nil {
		return nil, err
	}

	vboxHandle := resp.Returnval

	// Get Machine Handles
	machinesResp, err := service.IVirtualBox_getMachines(&IVirtualBox_getMachines{
		This: vboxHandle,
	})
	if err != nil {
		return nil, err
	}

	var results []vmModel

	for _, handle := range machinesResp.Returnval {
		// Get Name
		n, _ := service.IMachine_getName(&IMachine_getName{This: handle})
		
		// Get Memory (VirtualBox returns MB)
		m, _ := service.IMachine_getMemorySize(&IMachine_getMemorySize{This: handle})
		
		// Get CPU Count
		cp, _ := service.IMachine_getCPUCount(&IMachine_getCPUCount{This: handle})
		
		// Get State
		s, _ := service.IMachine_getState(&IMachine_getState{This: handle})

		d, _ := service.IMachine_getDescription(&IMachine_getDescription{This: handle})

		id, _ := service.IMachine_getId(&IMachine_getId{This: handle})

// func (service *vboxPortType) IMachine_getDescription(request *IMachine_getDescription) (*IMachine_getDescriptionResponse, error) {
// func (service *vboxPortType) IMachine_getId(request *IMachine_getId) (*IMachine_getIdResponse, error) {

// func (service *vboxPortType) IMachine_getOSTypeId(request *IMachine_getOSTypeId) (*IMachine_getOSTypeIdResponse, error) {
// func (service *vboxPortType) IMachine_getHardwareVersion(request *IMachine_getHardwareVersion) (*IMachine_getHardwareVersionResponse, error) {
// func (service *vboxPortType) IMachine_getHardwareUUID(request *IMachine_getHardwareUUID) (*IMachine_getHardwareUUIDResponse, error) {
// func (service *vboxPortType) IMachine_getCPUHotPlugEnabled(request *IMachine_getCPUHotPlugEnabled) (*IMachine_getCPUHotPlugEnabledResponse, error) {
// func (service *vboxPortType) IMachine_getCPUExecutionCap(request *IMachine_getCPUExecutionCap) (*IMachine_getCPUExecutionCapResponse, error) {
// func (service *vboxPortType) IMachine_getPointingHIDType(request *IMachine_getPointingHIDType) (*IMachine_getPointingHIDTypeResponse, error) {
// func (service *vboxPortType) IMachine_getKeyboardHIDType(request *IMachine_getKeyboardHIDType) (*IMachine_getKeyboardHIDTypeResponse, error) {
// func (service *vboxPortType) IMachine_getMediumAttachments(request *IMachine_getMediumAttachments) (*IMachine_getMediumAttachmentsResponse, error) {
// func (service *vboxPortType) IMachine_getUSBControllers(request *IMachine_getUSBControllers) (*IMachine_getUSBControllersResponse, error) {
// func (service *vboxPortType) IMachine_getStorageControllers(request *IMachine_getStorageControllers) (*IMachine_getStorageControllersResponse, error) {
// func (service *vboxPortType) IMachine_getSettingsFilePath(request *IMachine_getSettingsFilePath) (*IMachine_getSettingsFilePathResponse, error) {
// func (service *vboxPortType) IMachine_getSettingsAuxFilePath(request *IMachine_getSettingsAuxFilePath) (*IMachine_getSettingsAuxFilePathResponse, error) {
// func (service *vboxPortType) IMachine_getSettingsModified(request *IMachine_getSettingsModified) (*IMachine_getSettingsModifiedResponse, error) {
// func (service *vboxPortType) IMachine_getSessionState(request *IMachine_getSessionState) (*IMachine_getSessionStateResponse, error) {
// func (service *vboxPortType) IMachine_getSessionName(request *IMachine_getSessionName) (*IMachine_getSessionNameResponse, error) {
// func (service *vboxPortType) IMachine_getSessionPID(request *IMachine_getSessionPID) (*IMachine_getSessionPIDResponse, error) {
// func (service *vboxPortType) IMachine_getStateFilePath(request *IMachine_getStateFilePath) (*IMachine_getStateFilePathResponse, error) {
// func (service *vboxPortType) IMachine_getLogFolder(request *IMachine_getLogFolder) (*IMachine_getLogFolderResponse, error) {
// func (service *vboxPortType) IMachine_getSharedFolders(request *IMachine_getSharedFolders) (*IMachine_getSharedFoldersResponse, error) {
// func (service *vboxPortType) IMachine_getClipboardMode(request *IMachine_getClipboardMode) (*IMachine_getClipboardModeResponse, error) {
// func (service *vboxPortType) IMachine_getClipboardFileTransfersEnabled(request *IMachine_getClipboardFileTransfersEnabled) (*IMachine_getClipboardFileTransfersEnabledResponse, error) {


		var stateString string
		if s != nil && s.Returnval != nil {
			stateString = string(*s.Returnval)
		} else {
			stateString = "Unknown"
		}


		results = append(results, vmModel{
			ID:   types.StringValue(id.Returnval),
			Name:   types.StringValue(n.Returnval),
			Memory: types.Int64Value(int64(m.Returnval)),
			CPUs:   types.Int64Value(int64(cp.Returnval)),
			State:  types.StringValue(string(stateString)),
			Description: types.StringValue(d.Returnval),
		})
	}

	// Always Logoff!
	service.IWebsessionManager_logoff(&IWebsessionManager_logoff{
		RefIVirtualBox: vboxHandle,
	})

	return results, nil
}
package main

// gowsdl http://daddy.siwko.org:18083/?wsdl > vbox_interface.go
// on the virtualbox machine: vboxwebsrv -H 0.0.0.0 -A null

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hooklift/gowsdl/soap"
)

func WakeUpSubnet() {
	fmt.Fprintln(os.Stderr, "[ARP-WAKE] Initiating cross-platform subnet warming...")

	if runtime.GOOS == "linux" {
		// Linux specific optimization: Fire quick, non-blocking background shell pings
		// This forces the Linux bridge kernel module to track the target MAC neighbors safely
		cmdStr := `for ip in 50 51; do for host in {1..254}; do ping -c 1 -W 1 192.168.$ip.$host >/dev/null 2>&1 & done; done`
		cmd := exec.Command("bash", "-c", cmdStr)
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "[ARP-WARN] Linux sweep failed: %v\n", err)
		}
	} else {
		OldWakeUpSubnet()
	}

	// Essential pause to let the Linux interface finish updating the neighbor table states
	time.Sleep(4000 * time.Millisecond)
	fmt.Fprintln(os.Stderr, "[ARP-WAKE] Subnet warming phase completed.")
}
func OldWakeUpSubnet() {
	var wg sync.WaitGroup
	subnets := []string{"192.168.50", "192.168.51"}

	// fmt.Fprintln(os.Stderr, "[ARP-WAKE] Initiating throttled /23 subnet sweep...")

	for _, subnet := range subnets {
		for i := 1; i <= 254; i++ {
			wg.Add(1)
			go func(ipSuffix int, sub string) {
				defer wg.Done()

				targetIP := fmt.Sprintf("%s.%d", sub, ipSuffix)
				dst, err := net.ResolveUDPAddr("udp", targetIP+":7")
				if err != nil {
					return
				}

				conn, err := net.DialUDP("udp", nil, dst)
				if err != nil {
					return
				}
				defer conn.Close()

				_ = conn.SetDeadline(time.Now().Add(15 * time.Millisecond))
				_, _ = conn.Write([]byte{0})
			}(i, subnet)

			// Pacing delay: Prevents Windows socket buffer exhaustion
			time.Sleep(1 * time.Millisecond)
		}
	}

	wg.Wait()
	// Give the local switch a moment to catch the incoming hardware registrations
	time.Sleep(250 * time.Millisecond)
	// fmt.Fprintln(os.Stderr, "[ARP-WAKE] Subnet sweep completed successfully.")
}

// IPRegex safely extracts ipv4 sequences from lines
var ipRegex = regexp.MustCompile(`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`)

func lookupIPInARPCache(targetMac string) (string, bool) {
	if targetMac == "" {
		return "", false
	}

	// Normalize target MAC to plain hex characters (lowercase, no colons/hyphens)
	cleanTarget := strings.ToLower(strings.NewReplacer("-", "", ":", "").Replace(targetMac))

	// Execute standard platform ARP utility
	cmd := exec.Command("arp", "-a")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", false
	}

	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if line == "" {
			continue
		}

		// Normalize the entire row string to strip punctuation from potential MAC targets
		cleanLine := strings.NewReplacer("-", "", ":", "").Replace(line)

		// Check if this specific row contains our 12-character target hex MAC address
		if strings.Contains(cleanLine, cleanTarget) {
			// Extract the IP address safely out of the row line using regex
			foundIP := ipRegex.FindString(line)
			if foundIP != "" {
				return foundIP, true
			}
		}
	}

	return "", false
}
func OldlookupIPInARPCache(targetMac string) (string, bool) {
	if targetMac == "" {
		fmt.Fprintln(os.Stderr, "NO MAC SPECIFIED")
		return "", false
	}

	// Normalize target MAC to lowercase for safe string comparison
	targetMac = strings.ToLower(targetMac)

	// Execute the Windows standard ARP utility
	cmd := exec.Command("arp", "-a")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", false
	}

	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())

		// If this row contains our target hardware address...
		if strings.Contains(line, targetMac) {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				// Index 0: Internet Address (IP)
				// Index 1: Physical Address (MAC)
				ipAddress := fields[0]
				return ipAddress, true
			}
		}
	}

	return "", false
}
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

func formatMACAddress(rawMac string) string {
	if len(rawMac) != 12 {
		return rawMac
	}

	var parts []string
	for i := 0; i < 12; i += 2 {
		parts = append(parts, rawMac[i:i+2])
	}
	return strings.Join(parts, "-")
}

func (c *VBoxClient) GetDetailedVMs(ctx context.Context) ([]vmModel, error) {
	WakeUpSubnet()

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
		storage, _ := service.IMachine_getStorageControllers(&IMachine_getStorageControllers{This: handle})
		storageControllersList, _ := types.ListValueFrom(ctx, types.StringType, storage.Returnval)

		var disks []diskModel
		// 1. Get all Medium Attachments for this VM
		attachmentsResp, _ := service.IMachine_getMediumAttachments(&IMachine_getMediumAttachments{This: handle})

		for _, attachment := range attachmentsResp.Returnval {
			// 'attachment' contains the Slot, Port, and a handle to the Medium itself
			mediumHandle := attachment.Medium

			if mediumHandle != "" {
				locationResp, _ := service.IMedium_getLocation(&IMedium_getLocation{This: mediumHandle})
				sizeResp, _ := service.IMedium_getSize(&IMedium_getSize{This: mediumHandle})
				typeResp, _ := service.IMedium_getType(&IMedium_getType{This: mediumHandle})

				mediumType := "Unknown"
				if typeResp != nil && typeResp.Returnval != nil {
					mediumType = string(*typeResp.Returnval)
				}
				// Here you can create a diskModel and append it to a list of disks for the VM
				sizeInBytes := sizeResp.Returnval
				sizeInMB := sizeInBytes / (1024 * 1024)
				sizeInGB := sizeInBytes / (1024 * 1024 * 1024)
				disk := diskModel{
					Path:   locationResp.Returnval,
					Size:   sizeResp.Returnval,
					SizeMB: sizeInMB,
					SizeGB: sizeInGB,
					Type:   mediumType,
				}
				disks = append(disks, disk)
			}
		}
		disksList, _ := types.ListValueFrom(ctx, diskObjectType, disks)

		// --- READ GUEST IP ADDRESSES ---
		// --- DYNAMICALLY ENUMERATE GUEST PROPERTIES ---
		// --- FETCH VIRTUAL HARDWARE MAC ADDRESSES ---
		// --- FETCH VIRTUAL HARDWARE MAC ADDRESSES VIA HANDLES ---
		var rawIPs []string

		for slot := uint32(0); slot < 4; slot++ {
			adapterResp, err := service.IMachine_getNetworkAdapterContext(ctx, &IMachine_getNetworkAdapter{
				This: handle,
				Slot: slot,
			})

			if err == nil && adapterResp != nil && adapterResp.Returnval != "" {
				adapterHandle := adapterResp.Returnval

				// 1. Get the Attachment Type to see if the interface is plugged in
				attachResp, err := service.INetworkAdapter_getAttachmentTypeContext(ctx, &INetworkAdapter_getAttachmentType{
					This: adapterHandle,
				})

				// FIX: Check if the response pointer and Returnval pointer are safe to dereference
				if err == nil && attachResp != nil && attachResp.Returnval != nil {
					// Dereference the pointer and cast the custom type to a native string comparison
					attachmentStr := string(*attachResp.Returnval)

					if attachmentStr != "Null" {
						// 2. Fetch the flat hex MAC string from the adapter handle
						macResp, err := service.INetworkAdapter_getMACAddressContext(ctx, &INetworkAdapter_getMACAddress{
							This: adapterHandle,
						})

						// fmt.Fprintln(os.Stderr, "macResp.Returnval", macResp.Returnval)
						if err == nil && macResp != nil && macResp.Returnval != "" {
							// Transform "080027XXXXXX" into Windows ARP style "08-00-27-XX-XX-XX"
							formattedMac := formatMACAddress(macResp.Returnval)

							// 3. Query our local host routing cache for an IP pairing
							if ip, found := lookupIPInARPCache(formattedMac); found {
								rawIPs = append(rawIPs, ip)
							}
						}
					}
				}
			}
		}

		// --- FRAMEWORK-SAFE RESOLUTION BLOCK ---
		// FIX: Declaring and initializing ipsList clearly in the outer scope of the loop
		var ipsList types.List

		if len(rawIPs) == 0 {
			var diags diag.Diagnostics
			ipsList, diags = types.ListValue(types.StringType, []attr.Value{})
			if diags.HasError() {
				return nil, fmt.Errorf("failed to initialize empty list representation")
			}
		} else {
			var diags diag.Diagnostics
			ipsList, diags = types.ListValueFrom(ctx, types.StringType, rawIPs)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to process IP list mapping for VM %s", id.Returnval)
			}
		}
		// func (service *vboxPortType) IMachine_getDescription(request *IMachine_getDescription) (*IMachine_getDescriptionResponse, error) {
		// func (service *vboxPortType) IMachine_getId(request *IMachine_getId) (*IMachine_getIdResponse, error) {
		// func (service *vboxPortType) IMachine_getStorageControllers(request *IMachine_getStorageControllers) (*IMachine_getStorageControllersResponse, error) {

		// func (service *vboxPortType) IMachine_getOSTypeId(request *IMachine_getOSTypeId) (*IMachine_getOSTypeIdResponse, error) {
		// func (service *vboxPortType) IMachine_getHardwareVersion(request *IMachine_getHardwareVersion) (*IMachine_getHardwareVersionResponse, error) {
		// func (service *vboxPortType) IMachine_getHardwareUUID(request *IMachine_getHardwareUUID) (*IMachine_getHardwareUUIDResponse, error) {
		// func (service *vboxPortType) IMachine_getCPUHotPlugEnabled(request *IMachine_getCPUHotPlugEnabled) (*IMachine_getCPUHotPlugEnabledResponse, error) {
		// func (service *vboxPortType) IMachine_getCPUExecutionCap(request *IMachine_getCPUExecutionCap) (*IMachine_getCPUExecutionCapResponse, error) {
		// func (service *vboxPortType) IMachine_getPointingHIDType(request *IMachine_getPointingHIDType) (*IMachine_getPointingHIDTypeResponse, error) {
		// func (service *vboxPortType) IMachine_getKeyboardHIDType(request *IMachine_getKeyboardHIDType) (*IMachine_getKeyboardHIDTypeResponse, error) {
		// func (service *vboxPortType) IMachine_getMediumAttachments(request *IMachine_getMediumAttachments) (*IMachine_getMediumAttachmentsResponse, error) {
		// func (service *vboxPortType) IMachine_getUSBControllers(request *IMachine_getUSBControllers) (*IMachine_getUSBControllersResponse, error) {
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
			ID:                 types.StringValue(id.Returnval),
			Name:               types.StringValue(n.Returnval),
			Memory:             types.Int64Value(int64(m.Returnval)),
			CPUs:               types.Int64Value(int64(cp.Returnval)),
			State:              types.StringValue(string(stateString)),
			Description:        types.StringValue(d.Returnval),
			StorageControllers: storageControllersList,
			Disks:              disksList,
			IPAddresses:        ipsList,
		})
	}

	// Always Logoff!
	service.IWebsessionManager_logoff(&IWebsessionManager_logoff{
		RefIVirtualBox: vboxHandle,
	})
	return results, nil
}

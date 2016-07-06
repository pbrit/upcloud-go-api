package service

import (
	"github.com/jalle19/upcloud-go-sdk/upcloud"
	"github.com/jalle19/upcloud-go-sdk/upcloud/client"
	"github.com/jalle19/upcloud-go-sdk/upcloud/request"
	"os"
	"testing"
	"time"
)

// The service object used by the tests
var svc *Service

/**
Configures the test environment
*/
func init() {
	user, password := getCredentials()

	c := client.New(user, password)
	c.SetTimeout(time.Second * 30)
	svc = New(c)
}

/**
TestCreateModifyDeleteServer performs the following actions:

- creates a server
- modifies the server
- stops the server
- deletes the server

*/
func TestCreateModifyDeleteServer(t *testing.T) {
	t.Parallel()

	// Create a server
	serverDetails := createServer()
	t.Logf("Server %s with UUID %s created", serverDetails.Title, serverDetails.UUID)

	// Modify the server
	t.Log("Modifying the server ...")

	serverDetails, err := svc.ModifyServer(&request.ModifyServerRequest{
		UUID:  serverDetails.UUID,
		Title: "Modified server",
	})

	handleError(err)
	t.Log("Waiting for the server to exit maintenance state ...")

	err = svc.WaitForServerState(&request.WaitForServerStateRequest{
		UUID:         serverDetails.UUID,
		DesiredState: upcloud.ServerStateStarted,
		Timeout:      time.Minute * 5,
	})

	handleError(err)
	t.Logf("Server is now modified, new title is %s", serverDetails.Title)

	// Stop the server
	t.Log("Stopping the server ...")
	stopServer(serverDetails.UUID)
	t.Log("Server is now stopped")

	// Delete the server
	t.Log("Deleting the server ...")

	err = svc.DeleteServer(&request.DeleteServerRequest{
		UUID: serverDetails.UUID,
	})

	handleError(err)
	t.Log("Server is now deleted")
}

/**
TestCreateModifyDelete performs the following actions:

- creates a new storage disk
- modifies the storage
- deletes the storage

*/
func TestCreateModifyDelete(t *testing.T) {
	t.Parallel()

	// Create some storage
	storageDetails := createStorage()
	t.Logf("Storage %s with UUID %s created", storageDetails.Title, storageDetails.UUID)

	// Modify the storage
	t.Log("Modifying the storage ...")

	storageDetails, err := svc.ModifyStorage(&request.ModifyStorageRequest{
		UUID:  storageDetails.UUID,
		Title: "New fancy title",
	})

	handleError(err)
	t.Logf("Storage with UUID %s modified successfully, new title is %s", storageDetails.UUID, storageDetails.Title)

	// Delete the storage
	t.Log("Deleting the storage ...")

	err = svc.DeleteStorage(&request.DeleteStorageRequest{
		UUID: storageDetails.UUID,
	})

	handleError(err)
	t.Log("Storage is now deleted")
}

/**
TestAttachDetachStorage performs the following actions:

- creates a server
- stops the server
- creates a new storage disk
- attaches the storage
- detaches the storage
- deletes the storage
- deletes the server

*/
func TestAttachDetachStorage(t *testing.T) {
	t.Parallel()

	// Create a server
	serverDetails := createServer()
	t.Logf("Server %s with UUID %s created", serverDetails.Title, serverDetails.UUID)

	// Stop the server
	t.Log("Stopping the server ...")
	stopServer(serverDetails.UUID)
	t.Log("Server is now stopped")

	// Create some storage
	storageDetails := createStorage()
	t.Logf("Storage %s with UUID %s created", storageDetails.Title, storageDetails.UUID)

	// Attach the storage
	t.Logf("Attaching storage %s", storageDetails.UUID)

	serverDetails, err := svc.AttachStorageRequest(&request.AttachStorageRequest{
		StorageUUID: storageDetails.UUID,
		ServerUUID:  serverDetails.UUID,
		Type:        upcloud.StorageTypeDisk,
		Address:     "scsi:0:0",
	})

	handleError(err)
	t.Logf("Storage attached to server with UUID %s", serverDetails.UUID)

	// Detach the storage
	t.Logf("Detaching storage %s", storageDetails.UUID)

	serverDetails, err = svc.DetachStorageRequest(&request.DetachStorageRequest{
		ServerUUID: serverDetails.UUID,
		Address:    "scsi:0:0",
	})

	handleError(err)
	t.Logf("Storage %s detached", storageDetails.UUID)

	// Delete the storage
	t.Log("Deleting the storage ...")

	err = svc.DeleteStorage(&request.DeleteStorageRequest{
		UUID: storageDetails.UUID,
	})

	handleError(err)
	t.Log("Storage is now deleted")

	// Delete the server
	t.Log("Deleting the server ...")

	err = svc.DeleteServer(&request.DeleteServerRequest{
		UUID: serverDetails.UUID,
	})

	handleError(err)
	t.Log("Server is now deleted")
}

/**
TestLoadEjectCDROM performs the following actions:

- creates a server
- stops the server
- attaches a CD-ROM device
- starts the server
- loads a CD-ROM
- ejects the CD-ROM
- stops the server
- deletes the server

*/
func TestLoadEjectCDROM(t *testing.T) {
	t.Parallel()

	// Create the server
	serverDetails := createServer()
	t.Logf("Server %s with UUID %s created", serverDetails.Title, serverDetails.UUID)

	// Stop the server
	t.Log("Stopping the server ...")
	stopServer(serverDetails.UUID)
	t.Log("Server is now stopped")

	// Attach CD-ROM device
	t.Logf("Attaching CD-ROM device to server with UUID %s", serverDetails.UUID)
	serverDetails, err := svc.AttachStorageRequest(&request.AttachStorageRequest{
		ServerUUID: serverDetails.UUID,
		Type:       upcloud.StorageTypeCDROM,
	})

	handleError(err)
	t.Log("CD-ROM is now attached")

	// Start the server
	t.Logf("Starting server with UUID %s", serverDetails.UUID)
	serverDetails, err = svc.StartServer(&request.StartServerRequest{
		UUID:    serverDetails.UUID,
		Timeout: time.Minute * 5,
	})

	handleError(err)

	err = svc.WaitForServerState(&request.WaitForServerStateRequest{
		UUID:         serverDetails.UUID,
		DesiredState: upcloud.ServerStateStarted,
		Timeout:      time.Minute * 5,
	})

	handleError(err)
	t.Log("Server is now started")

	// Load the CD-ROM
	t.Log("Loading CD-ROM into CD-ROM device")
	serverDetails, err = svc.LoadCDROM(&request.LoadCDROMRequest{
		ServerUUID:  serverDetails.UUID,
		StorageUUID: "01000000-0000-4000-8000-000030060101",
	})

	handleError(err)
	t.Log("CD-ROM is now loaded")

	// Eject the CD-ROM
	t.Log("Ejecting CD-ROM from CD-ROM device")
	serverDetails, err = svc.EjectCDROM(&request.EjectCDROMRequest{
		ServerUUID: serverDetails.UUID,
	})

	handleError(err)
	t.Log("CD-ROM is now ejected")

	// Stop the server
	t.Logf("Stopping server with UUID %s", serverDetails.UUID)
	stopServer(serverDetails.UUID)
	t.Log("Server is now stopped")

	// Delete the server
	// Delete the server
	t.Log("Deleting the server ...")

	err = svc.DeleteServer(&request.DeleteServerRequest{
		UUID: serverDetails.UUID,
	})

	handleError(err)
	t.Log("Server is now deleted")
}

/**
Creates a server and returns the details about it, panic if creation fails
*/
func createServer() *upcloud.ServerDetails {
	createServerRequest := request.CreateServerRequest{
		Zone:             "fi-hel1",
		Title:            "Integration test server #1",
		Hostname:         "debian.example.com",
		PasswordDelivery: request.PasswordDeliveryNone,
		StorageDevices: []request.CreateServerStorageDevice{
			{
				Action:  request.CreateStorageDeviceActionClone,
				Storage: "01000000-0000-4000-8000-000030060200",
				Title:   "disk1",
				Size:    30,
				Tier:    request.CreateStorageDeviceTierMaxIOPS,
			},
		},
		IPAddresses: []request.CreateServerIPAddress{
			{
				Access: upcloud.IPAddressAccessPrivate,
				Family: upcloud.IPAddressFamilyIPv4,
			},
			{
				Access: upcloud.IPAddressAccessPublic,
				Family: upcloud.IPAddressFamilyIPv4,
			},
			{
				Access: upcloud.IPAddressAccessPublic,
				Family: upcloud.IPAddressFamilyIPv6,
			},
		},
	}

	// Create the server and block until it has started
	serverDetails, err := svc.CreateServer(&createServerRequest)

	if err != nil {
		panic(err)
	}

	// Wait for the server to start
	err = svc.WaitForServerState(&request.WaitForServerStateRequest{
		UUID:         serverDetails.UUID,
		DesiredState: upcloud.ServerStateStarted,
		Timeout:      time.Minute * 5,
	})

	handleError(err)

	return serverDetails
}

/**
Stops the specified server
*/
func stopServer(uuid string) {
	serverDetails, err := svc.StopServer(&request.StopServerRequest{
		UUID:    uuid,
		Timeout: time.Minute * 5,
	})

	handleError(err)

	err = svc.WaitForServerState(&request.WaitForServerStateRequest{
		UUID:         serverDetails.UUID,
		DesiredState: upcloud.ServerStateStopped,
		Timeout:      time.Minute * 5,
	})

	handleError(err)
}

/**
Creates a piece of storage and returns the details about it, panic if creation fails
*/
func createStorage() *upcloud.StorageDetails {
	createStorageRequest := request.CreateStorageRequest{
		Tier:  upcloud.StorageTierMaxIOPS,
		Title: "Test storage",
		Size:  50,
		Zone:  "fi-hel1",
	}

	storageDetails, err := svc.CreateStorage(&createStorageRequest)

	if err != nil {
		panic(err)
	}

	return storageDetails
}

/**
Handles the error by panicing, thus stopping the test execution
*/
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

/**
Reads the API username and password from the environment, panics if they are not available
*/
func getCredentials() (string, string) {
	user := os.Getenv("UPCLOUD_GO_SDK_TEST_USER")
	password := os.Getenv("UPCLOUD_GO_SDK_TEST_PASSWORD")

	if user == "" || password == "" {
		panic("Unable to retrieve credentials from the environment, ensure UPCLOUD_GO_SDK_TEST_USER and UPCLOUD_GO_SDK_TEST_PASSWORD are exported")
	}

	return user, password
}

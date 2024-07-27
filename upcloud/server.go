package upcloud

import (
	"encoding/json"
)

// Constants
const (
	ServerStateStarted     = "started"
	ServerStateStopped     = "stopped"
	ServerStateMaintenance = "maintenance"
	ServerStateError       = "error"

	VideoModelVGA    = "vga"
	VideoModelCirrus = "cirrus"

	NICModelE1000   = "e1000"
	NICModelVirtio  = "virtio"
	NICModelRTL8139 = "rtl8139"

	StopTypeSoft = "soft"
	StopTypeHard = "hard"

	RemoteAccessTypeVNC   = "vnc"
	RemoteAccessTypeSPICE = "spice"
)

// ServerConfigurations represents a /server_size response
type ServerConfigurations struct {
	ServerConfigurations []ServerConfiguration `json:"server_sizes"`
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *ServerConfigurations) UnmarshalJSON(b []byte) error {
	type serverConfigurationWrapper struct {
		ServerConfigurations []ServerConfiguration `json:"server_size"`
	}

	v := struct {
		ServerConfigurations serverConfigurationWrapper `json:"server_sizes"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	s.ServerConfigurations = v.ServerConfigurations.ServerConfigurations

	return nil
}

// ServerConfiguration represents a server configuration
type ServerConfiguration struct {
	CoreNumber   int `json:"core_number,string"`
	MemoryAmount int `json:"memory_amount,string"`
}

// Servers represents a /server response
type Servers struct {
	Servers []Server `json:"servers"`
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *Servers) UnmarshalJSON(b []byte) error {
	type serverWrapper struct {
		Servers []Server `json:"server"`
	}

	v := struct {
		Servers serverWrapper `json:"servers"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	s.Servers = v.Servers.Servers

	return nil
}

// ServerTagSlice is a slice of string.
// It exists to allow for a custom JSON unmarshaller.
type ServerTagSlice []string

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (t *ServerTagSlice) UnmarshalJSON(b []byte) error {
	v := struct {
		Tags []string `json:"tag"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*t) = v.Tags

	return nil
}

// Server represents a server

// +k8s:deepcopy-gen=true

type Server struct {
	CoreNumber   string         `json:"core_number,omitempty"`
	Hostname     string         `json:"hostname,omitempty"`
	License      float64        `json:"license,omitempty"`
	MemoryAmount string         `json:"memory_amount,omitempty"`
	Plan         string         `json:"plan,omitempty"`
	Progress     string         `json:"progress,omitempty"`
	State        string         `json:"state,omitempty"`
	Tags         ServerTagSlice `json:"tags,omitempty"`
	Title        string         `json:"title,omitempty"`
	UUID         string         `json:"uuid,omitempty"`
	Zone         string         `json:"zone,omitempty"`
}

// ServerStorageDeviceSlice is a slice of ServerStorageDevices.
// It exists to allow for a custom JSON unmarshaller.
type ServerStorageDeviceSlice []ServerStorageDevice

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *ServerStorageDeviceSlice) UnmarshalJSON(b []byte) error {
	v := struct {
		StorageDevices []ServerStorageDevice `json:"storage_device"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*s) = v.StorageDevices

	return nil
}

// ServerNetworking represents the networking on a server response.
// It is castable to a Networking struct.

// +k8s:deepcopy-gen=true

type ServerNetworking Networking

// ServerDetails represents details about a server

// +k8s:deepcopy-gen=true

type ServerDetails struct {
	Server `json:",inline,omitempty"`

	BootOrder string `json:"boot_order,omitempty"`
	// TODO: Convert to boolean
	Firewall             string                   `json:"firewall,omitempty"`
	Host                 int                      `json:"host,omitempty"`
	IPAddresses          IPAddressSlice           `json:"ip_addresses,omitempty"`
	Labels               LabelSlice               `json:"labels,omitempty"`
	Metadata             Boolean                  `json:"metadata,omitempty"`
	NICModel             string                   `json:"nic_model,omitempty"`
	Networking           ServerNetworking         `json:"networking,omitempty"`
	ServerGroup          string                   `json:"server_group,omitempty"`
	SimpleBackup         string                   `json:"simple_backup,omitempty"`
	StorageDevices       ServerStorageDeviceSlice `json:"storage_devices,omitempty"`
	Timezone             string                   `json:"timezone,omitempty"`
	VideoModel           string                   `json:"video_model,omitempty"`
	RemoteAccessEnabled  Boolean                  `json:"remote_access_enabled,omitempty"`
	RemoteAccessType     string                   `json:"remote_access_type,omitempty"`
	RemoteAccessHost     string                   `json:"remote_access_host,omitempty"`
	RemoteAccessPassword string                   `json:"remote_access_password,omitempty"`
	RemoteAccessPort     string                   `json:"remote_access_port,omitempty"`
}

func (s *ServerDetails) StorageDevice(storageUUID string) *ServerStorageDevice {
	for _, storageDevice := range s.StorageDevices {
		if storageDevice.UUID == storageUUID {
			return &storageDevice
		}
	}
	return nil
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *ServerDetails) UnmarshalJSON(b []byte) error {
	type localServerDetails ServerDetails

	v := struct {
		ServerDetails localServerDetails `json:"server"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*s) = ServerDetails(v.ServerDetails)

	return nil
}

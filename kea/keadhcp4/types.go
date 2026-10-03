package keadhcp4

import (
	"encoding/json"
	"net/netip"

	"github.com/venthe/terraform-provider-kea/kea"
)

type Subnet struct {
	ID     uint32 `json:"id"`
	Subnet string `json:"subnet"`

	SharedNetworkName string `json:"shared-network-name,omitempty"`
}

type Reservation struct {
	SubnetID uint32 `json:"subnet-id"`

	CircuitID kea.HexID `json:"circuit-id,omitempty"`
	ClientID  kea.HexID `json:"client-id,omitempty"`
	DUID      kea.HexID `json:"duid,omitempty"`
	FlexID    kea.HexID `json:"flex-id,omitempty"`
	HWAddress kea.HexID `json:"hw-address,omitempty"`

	BootFileName   *string          `json:"boot-file-name,omitempty"`
	ClientClasses  []string         `json:"client-classes,omitempty"`
	Hostname       *string          `json:"hostname,omitempty"`
	IPAddress      *netip.Addr      `json:"ip-address,omitempty"`
	NextServer     *netip.Addr      `json:"next-server,omitempty"`
	OptionData     []OptionData     `json:"option-data,omitempty"`
	ServerHostname *string          `json:"server-hostname,omitempty"`
	UserContext    *json.RawMessage `json:"user-context,omitempty"`
}

type OptionData struct {
	Name          *string  `json:"name,omitempty"`
	Data          *string  `json:"data,omitempty"`
	Code          *uint8   `json:"code,omitempty"`
	Space         *string  `json:"space,omitempty"`
	CSVFormat     *bool    `json:"csv-format,omitempty"`
	AlwaysSend    *bool    `json:"always-send,omitempty"`
	NeverSend     *bool    `json:"never-send,omitempty"`
	ClientClasses []string `json:"client-classes,omitempty"`
}

func (r *Reservation) UnmarshalJSON(data []byte) error {
	type Alias Reservation
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// The Kea API returns these fields as empty strings when not set,
	// rather than omitting them from the response.
	r.BootFileName = nilIfEmpty(r.BootFileName)
	r.Hostname = nilIfEmpty(r.Hostname)
	r.ServerHostname = nilIfEmpty(r.ServerHostname)

	// The Kea API returns "0.0.0.0" for next-server when not set.
	if r.NextServer != nil && (!r.NextServer.IsValid() || r.NextServer.IsUnspecified()) {
		r.NextServer = nil
	}

	return nil
}

func nilIfEmpty(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}

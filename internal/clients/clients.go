package clients

import "github.com/venthe/terraform-provider-kea/kea/keadhcp4"

type KeaClients struct {
	DHCP4 *keadhcp4.Client
}

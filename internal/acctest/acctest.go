package acctest

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/venthe/terraform-provider-kea/internal/provider"
)

// ProtoV6ProviderFactories is used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"kea": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// PreCheck validates required environment variables before running acceptance tests.
// Assumes the Kea server is configured with testing/kea-dhcp4.conf.
func PreCheck(t *testing.T) {
	t.Helper()

	address := os.Getenv("KEA_DHCP4_ADDRESS")
	if address == "" {
		t.Fatal("KEA_DHCP4_ADDRESS must be set for acceptance tests")
	}
}

func ProviderConfig() string {
	return `
provider "kea" {}
`
}

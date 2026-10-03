package dhcp4_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/venthe/terraform-provider-kea/internal/acctest"
	"github.com/venthe/terraform-provider-kea/kea"
	"github.com/venthe/terraform-provider-kea/kea/keadhcp4"
	"github.com/venthe/terraform-provider-kea/kea/keaquery"
)

func TestAccReservation_basic(t *testing.T) {
	mac := "02:A3:7B:4E:91:22"
	ip := "10.67.0.42"
	resourceName := "kea_dhcp4_reservation.test"
	query := keaquery.ReservationByIdentifier(1, "hw-address", mac)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_basic(1, mac, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "hw_address", mac),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
				PostApplyFunc: testAccCheckReservationExists(t, query),
			},
		},
	})
}

func TestAccReservation_withClientID(t *testing.T) {
	clientID := "01:aa:bb:cc"
	ip := "10.67.0.74"
	resourceName := "kea_dhcp4_reservation.test"
	query := keaquery.ReservationByIdentifier(1, "client-id", clientID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_withClientID(1, clientID, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "client_id", clientID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
				PostApplyFunc: testAccCheckReservationExists(t, query),
			},
		},
	})
}

func TestAccReservation_withCircuitID(t *testing.T) {
	circuitID := "01:02:03:04"
	ip := "10.67.0.73"
	resourceName := "kea_dhcp4_reservation.test"
	query := keaquery.ReservationByIdentifier(1, "circuit-id", circuitID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_withCircuitID(1, circuitID, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "circuit_id", circuitID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
				PostApplyFunc: testAccCheckReservationExists(t, query),
			},
		},
	})
}

func TestAccReservation_withDUID(t *testing.T) {
	duid := "00:03:00:01:de:ad:be:ef:ca:fe"
	ip := "10.67.0.75"
	resourceName := "kea_dhcp4_reservation.test"
	query := keaquery.ReservationByIdentifier(1, "duid", duid)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_withDUID(1, duid, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "duid", duid),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
				PostApplyFunc: testAccCheckReservationExists(t, query),
			},
		},
	})
}

func TestAccReservation_withFlexID(t *testing.T) {
	flexID := "01:02:03:04:05:06"
	ip := "10.67.0.76"
	resourceName := "kea_dhcp4_reservation.test"
	query := keaquery.ReservationByIdentifier(1, "flex-id", flexID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_withFlexID(1, flexID, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "flex_id", flexID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
				PostApplyFunc: testAccCheckReservationExists(t, query),
			},
		},
	})
}

func TestAccReservation_withHostname(t *testing.T) {
	mac := "02:f8:c2:5d:19:a6"
	ip := "10.67.0.142"
	hostname := fmt.Sprintf("test-host-%d", rand.Intn(10000))
	resourceName := "kea_dhcp4_reservation.test"
	query := keaquery.ReservationByIdentifier(1, "hw-address", mac)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_withHostname(1, mac, ip, hostname),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "hw_address", mac),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttr(resourceName, "hostname", hostname),
				),
				PostApplyFunc: testAccCheckReservationExists(t, query),
			},
		},
	})
}

func TestAccReservation_withOptionData(t *testing.T) {
	mac := "02:6b:d9:31:84:cf"
	ip := "10.67.0.27"
	resourceName := "kea_dhcp4_reservation.test"
	query := keaquery.ReservationByIdentifier(1, "hw-address", mac)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_withOptionData(1, mac, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "hw_address", mac),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttr(resourceName, "option_data.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "option_data.0.name", "domain-name-servers"),
				),
				PostApplyFunc: testAccCheckReservationExists(t, query),
			},
		},
	})
}

func TestAccReservation_destroy(t *testing.T) {
	mac := "02:9a:3f:6c:d1:84"
	ip := "10.67.0.201"
	query := keaquery.ReservationByIdentifier(1, "hw-address", mac)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_basic(1, mac, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("kea_dhcp4_reservation.test", "hw_address", mac),
				),
				PostApplyFunc: testAccCheckReservationExists(t, query),
			},
			{
				Config:        acctest.ProviderConfig(),
				PostApplyFunc: testAccCheckReservationDestroyed(t, query),
			},
		},
	})
}

func TestAccReservation_disappears(t *testing.T) {
	mac := "02:1c:8e:a7:f2:3d"
	ip := "10.67.0.177"
	resourceName := "kea_dhcp4_reservation.test"
	query := keaquery.ReservationByIdentifier(1, "hw-address", mac)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_basic(1, mac, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "hw_address", mac),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
				),
				PostApplyFunc: testAccDeleteReservation(t, query),
			},
			{
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccReservation_global(t *testing.T) {
	mac := "02:95:4a:62:b8:e1"
	ip := "192.168.67.143"
	resourceName := "kea_dhcp4_reservation.test"
	query := keaquery.ReservationByIdentifier(0, "hw-address", mac)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_global(mac, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "subnet_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "hw_address", mac),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
				PostApplyFunc: testAccCheckReservationExists(t, query),
			},
		},
	})
}

func TestAccReservation_update(t *testing.T) {
	mac := "02:d4:71:39:ac:56"
	ip := "10.67.0.91"
	hostname1 := fmt.Sprintf("test-host-%d", rand.Intn(10000))
	hostname2 := fmt.Sprintf("test-host-%d", rand.Intn(10000))
	resourceName := "kea_dhcp4_reservation.test"
	query := keaquery.ReservationByIdentifier(1, "hw-address", mac)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_withHostname(1, mac, ip, hostname1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "hostname", hostname1),
				),
				PostApplyFunc: testAccCheckReservationExists(t, query),
			},
			{
				Config: testAccReservationConfig_withHostname(1, mac, ip, hostname2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "hostname", hostname2),
				),
				PostApplyFunc: testAccCheckReservationExists(t, query),
			},
		},
	})
}

func TestAccReservation_reorderSetsNoUpdate(t *testing.T) {
	mac := "02:e7:2f:58:c3:9b"
	ip := "10.67.0.165"
	resourceName := "kea_dhcp4_reservation.test"
	query := keaquery.ReservationByIdentifier(1, "hw-address", mac)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_withClientClassesAndOptions(1, mac, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "hw_address", mac),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttr(resourceName, "client_classes.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "option_data.#", "3"),
				),
				PostApplyFunc: testAccCheckReservationExists(t, query),
			},
			{
				Config: testAccReservationConfig_withClientClassesAndOptionsReordered(1, mac, ip),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccReservation_withUserContext(t *testing.T) {
	mac := "02:5a:b6:4d:e9:72"
	ip := "10.67.0.64"
	resourceName := "kea_dhcp4_reservation.test"
	query := keaquery.ReservationByIdentifier(1, "hw-address", mac)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReservationConfig_withUserContext(1, mac, ip),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "hw_address", mac),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ip),
					resource.TestCheckResourceAttrSet(resourceName, "user_context"),
				),
				PostApplyFunc: testAccCheckReservationExists(t, query),
			},
		},
	})
}

func testAccKeaClient() *keadhcp4.Client {
	transport := &kea.HTTPTransport{
		Endpoint: os.Getenv("KEA_DHCP4_ADDRESS"),
		Username: os.Getenv("KEA_DHCP4_HTTP_USERNAME"),
		Password: os.Getenv("KEA_DHCP4_HTTP_PASSWORD"),
	}
	return keadhcp4.NewClient(transport)
}

func testAccCheckReservationExists(t *testing.T, query keaquery.ReservationQuery) func() {
	return func() {
		client := testAccKeaClient()

		reservation, err := client.GetReservation(context.Background(), query)
		if err != nil {
			t.Fatalf("failed to get reservation: %v", err)
		}

		if reservation == nil {
			t.Fatal("reservation does not exist in Kea")
		}
	}
}

func testAccCheckReservationDestroyed(t *testing.T, query keaquery.ReservationQuery) func() {
	return func() {
		client := testAccKeaClient()

		reservation, err := client.GetReservation(context.Background(), query)
		if err != nil {
			t.Fatalf("failed to check reservation: %v", err)
		}

		if reservation != nil {
			t.Fatal("reservation still exists")
		}
	}
}

func testAccDeleteReservation(t *testing.T, query keaquery.ReservationQuery) func() {
	return func() {
		client := testAccKeaClient()
		if err := client.DeleteReservation(context.Background(), query); err != nil {
			t.Fatalf("failed to delete reservation: %v", err)
		}
	}
}

func testAccReservationConfig_basic(subnetID uint32, mac, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id  = %d
  hw_address = %q
  ip_address = %q
}
`, acctest.ProviderConfig(), subnetID, mac, ip)
}

func testAccReservationConfig_withClientID(subnetID uint32, clientID, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id  = %d
  client_id  = %q
  ip_address = %q
}
`, acctest.ProviderConfig(), subnetID, clientID, ip)
}

func testAccReservationConfig_withCircuitID(subnetID uint32, circuitID, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id  = %d
  circuit_id = %q
  ip_address = %q
}
`, acctest.ProviderConfig(), subnetID, circuitID, ip)
}

func testAccReservationConfig_withDUID(subnetID uint32, duid, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id  = %d
  duid       = %q
  ip_address = %q
}
`, acctest.ProviderConfig(), subnetID, duid, ip)
}

func testAccReservationConfig_withFlexID(subnetID uint32, flexID, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id  = %d
  flex_id    = %q
  ip_address = %q
}
`, acctest.ProviderConfig(), subnetID, flexID, ip)
}

func testAccReservationConfig_withHostname(subnetID uint32, mac, ip, hostname string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id  = %d
  hw_address = %q
  ip_address = %q
  hostname   = %q
}
`, acctest.ProviderConfig(), subnetID, mac, ip, hostname)
}

func testAccReservationConfig_withOptionData(subnetID uint32, mac, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id  = %d
  hw_address = %q
  ip_address = %q

  option_data {
    name       = "domain-name-servers"
    data       = "8.8.8.8,8.8.4.4"
    csv_format = true
  }
}
`, acctest.ProviderConfig(), subnetID, mac, ip)
}

func testAccReservationConfig_global(mac, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id  = 0
  hw_address = %q
  ip_address = %q
}
`, acctest.ProviderConfig(), mac, ip)
}

func testAccReservationConfig_withClientClassesAndOptions(subnetID uint32, mac, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id      = %d
  hw_address     = %q
  ip_address     = %q
  client_classes = ["web-servers", "production", "monitoring"]

  option_data {
    name = "domain-name-servers"
    data = "8.8.8.8,8.8.4.4"
  }

  option_data {
    name = "domain-name"
    data = "example.com"
  }

  option_data {
    name = "routers"
    data = "192.0.2.1"
  }
}
`, acctest.ProviderConfig(), subnetID, mac, ip)
}

func testAccReservationConfig_withClientClassesAndOptionsReordered(subnetID uint32, mac, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id      = %d
  hw_address     = %q
  ip_address     = %q
  client_classes = ["monitoring", "production", "web-servers"]

  option_data {
    name = "routers"
    data = "192.0.2.1"
  }

  option_data {
    name = "domain-name-servers"
    data = "8.8.8.8,8.8.4.4"
  }

  option_data {
    name = "domain-name"
    data = "example.com"
  }
}
`, acctest.ProviderConfig(), subnetID, mac, ip)
}

func testAccReservationConfig_withUserContext(subnetID uint32, mac, ip string) string {
	return fmt.Sprintf(`
%s

resource "kea_dhcp4_reservation" "test" {
  subnet_id    = %d
  hw_address   = %q
  ip_address   = %q
  user_context = jsonencode({
    department = "engineering"
    owner      = "network-team"
    tags       = ["production", "critical"]
  })
}
`, acctest.ProviderConfig(), subnetID, mac, ip)
}

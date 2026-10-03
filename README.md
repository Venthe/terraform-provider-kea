# Terraform Provider for Kea DHCP

A Terraform provider for managing [Kea DHCP](https://www.isc.org/kea/) server
resources.

## Overview

This provider enables Infrastructure as Code management of Kea DHCP servers
through their control channel API. It allows you to manage DHCP reservations
programmatically without requiring server restarts.

Currently supported:
- DHCPv4 host reservations (resource and data source)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24 (for building from source)
- A Kea DHCP server with a control channel configured

### Kea Server Configuration

This provider communicates directly with Kea DHCP servers via their control
channel. You must configure either:

- **UNIX domain socket** - For local connections to the server
- **HTTP control socket** - For remote connections (requires the `http` hook
  library)

**Important:** This provider does **not** work with the Kea Control Agent
(`kea-ctrl-agent`). The Control Agent is deprecated in favor of direct HTTP
control sockets on the DHCP services themselves. Configure the `http` hook
library directly on your `kea-dhcp4` (or `kea-dhcp6`) service instead.

See [`testing/kea-dhcp4.conf`](testing/kea-dhcp4.conf) for a working example of
control socket and hook library configuration.

## Installation

### Terraform Registry

```terraform
terraform {
  required_providers {
    kea = {
      source = "venthe/kea"
    }
  }
}
```

### Building from Source

```sh
git clone https://github.com/venthe/terraform-provider-kea.git
cd terraform-provider-kea
make build
```

To use a locally built provider, add a dev override to your `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "venthe/kea" = "/path/to/terraform-provider-kea"
  }
  direct {}
}
```

## Quick Start

Configure the provider to connect to your Kea server via HTTP:

```terraform
provider "kea" {
  dhcp4 = {
    address = "http://kea-server:8000"
  }
}
```

Or via UNIX socket for local connections:

```terraform
provider "kea" {
  dhcp4 = {
    address = "unix:///run/kea/kea-dhcp4.sock"
  }
}
```

Or via environment variables:

```sh
export KEA_DHCP4_ADDRESS="http://kea-server:8000"
export KEA_DHCP4_HTTP_USERNAME="admin"  # optional
export KEA_DHCP4_HTTP_PASSWORD="secret" # optional
```

```terraform
provider "kea" {}
```

### HTTPS and mutual TLS

For `https://` addresses, the TLS options below are available. Each falls back to an environment variable.

| Attribute              | Environment variable              |
| ---------------------- | --------------------------------- |
| `tls_ca_file`          | `KEA_DHCP4_TLS_CA_FILE`           |
| `tls_client_cert_file` | `KEA_DHCP4_TLS_CLIENT_CERT_FILE`  |
| `tls_client_key_file`  | `KEA_DHCP4_TLS_CLIENT_KEY_FILE`   |
| `tls_server_name`      | `KEA_DHCP4_TLS_SERVER_NAME`       |

`tls_client_cert_file` and `tls_client_key_file` must be set together.

### Passing values via tfvars

Declare variables and reference them in the provider block:

```terraform
variable "kea_address"         { type = string }
variable "kea_tls_ca_file"     { type = string }
variable "kea_tls_client_cert" { type = string }
variable "kea_tls_client_key"  { type = string }
variable "kea_tls_server_name" { type = string }

provider "kea" {
  dhcp4 = {
    address              = var.kea_address
    tls_ca_file          = var.kea_tls_ca_file
    tls_client_cert_file = var.kea_tls_client_cert
    tls_client_key_file  = var.kea_tls_client_key
    tls_server_name      = var.kea_tls_server_name
  }
}
```

Set them in `terraform.tfvars` or `*.auto.tfvars` (a file named `.auto.tfvar` is NOT loaded automatically):

```hcl
kea_address         = "https://172.22.0.1:8000"
kea_tls_ca_file     = "./issuer-ca.crt"
kea_tls_client_cert = "./docker.chain.crt"
kea_tls_client_key  = "./docker.key"
kea_tls_server_name = "edge01.home.arpa"
```

Or use `TF_VAR_kea_tls_client_key=...` environment variables.

Create a DHCP reservation:

```terraform
resource "kea_dhcp4_reservation" "webserver" {
  subnet_id  = 1
  hw_address = "aa:bb:cc:dd:ee:01"
  ip_address = "192.0.2.10"
  hostname   = "webserver"
}
```

See the [`examples/`](examples/) directory for more usage examples.

## Documentation

Full documentation is available in the [`docs/`](docs/) directory.

## License

Mozilla Public License 2.0 - see [LICENSE](LICENSE) for details.

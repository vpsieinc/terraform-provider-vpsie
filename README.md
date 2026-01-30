# Terraform Provider for VPSie

Terraform provider for managing [VPSie](https://vpsie.com/) cloud infrastructure resources.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24 (to build the provider)
- A VPSie account and API access token

## Setup

### Authentication

The provider requires a VPSie API access token. You can configure it in two ways:

**Option 1: Provider configuration**

```hcl
provider "vpsie" {
  access_token = var.vpsie_access_token
}
```

**Option 2: Environment variable**

```bash
export VPSIE_ACCESS_TOKEN="your-access-token"
```

### Provider Configuration

```hcl
terraform {
  required_providers {
    vpsie = {
      source = "vpsie/vpsie"
    }
  }
}

provider "vpsie" {}
```

## Usage Example

```hcl
# Create a server
resource "vpsie_server" "example" {
  hostname            = "my-server"
  dc_identifier       = "dc-identifier"
  os_identifier       = "os-identifier"
  resource_identifier = "resource-identifier"
  project_id          = 1
  password            = "secure-password"
  delete_reason       = "no longer needed"
}

# Create a storage volume
resource "vpsie_storage" "example" {
  name = "my-volume"
  size = 20
}
```

## Supported Resources

### Resources

| Resource | Description |
|----------|-------------|
| `vpsie_server` | Compute instances |
| `vpsie_storage` | Block storage volumes |
| `vpsie_storage_attachment` | Attach storage to servers |
| `vpsie_storage_snapshot` | Storage volume snapshots |
| `vpsie_vpc` | Virtual private clouds |
| `vpsie_vpc_server_assignment` | Assign servers to VPCs |
| `vpsie_firewall` | Firewall rule groups |
| `vpsie_firewall_attachment` | Attach firewalls to servers |
| `vpsie_kubernetes` | Kubernetes clusters |
| `vpsie_kubernetes_group` | Kubernetes node groups |
| `vpsie_loadbalancer` | Load balancers |
| `vpsie_domain` | DNS domains |
| `vpsie_dns_record` | DNS records |
| `vpsie_reverse_dns` | Reverse DNS entries |
| `vpsie_gateway` | Network gateways |
| `vpsie_floating_ip` | Floating IP addresses |
| `vpsie_image` | Custom images |
| `vpsie_snapshot` | Server snapshots |
| `vpsie_snapshot_policy` | Snapshot scheduling policies |
| `vpsie_backup` | Server backups |
| `vpsie_backup_policy` | Backup scheduling policies |
| `vpsie_sshkey` | SSH keys |
| `vpsie_script` | Startup scripts |
| `vpsie_project` | Projects |
| `vpsie_bucket` | Object storage buckets |
| `vpsie_monitoring_rule` | Monitoring alert rules |
| `vpsie_access_token` | API access tokens |

### Data Sources

All resources above have corresponding data sources (e.g., `data.vpsie_server`) for reading existing infrastructure. Additional read-only data sources:

- `data.vpsie_datacenter` — Available data center locations
- `data.vpsie_ip` — IP address information

## Development

### Building

```bash
go build -v .
```

### Running Tests

Acceptance tests require a VPSie API token and will create real resources:

```bash
export VPSIE_ACCESS_TOKEN="your-access-token"
TF_ACC=1 go test ./... -v -timeout 120m
```

Run a single test:

```bash
TF_ACC=1 go test ./internal/services/storage -run TestAccStorageResource -v -timeout 120m
```

### Generating Documentation

```bash
go generate ./...
```

This formats example HCL files and regenerates the `docs/` directory from schemas.

### Linting

```bash
golangci-lint run
```

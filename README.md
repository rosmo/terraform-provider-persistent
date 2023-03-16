# Persistent Counter provider for Terraform

```
terraform {
  required_providers {
    persistent = {
      source = "rosmo/persistent"
    }
  }
}

resource "persistent_counter" "example" {
  keys     = ["a", "b", "c"]
}

# Result would be:
#   persistent_counter.example.values = { a = 0, b = 1, c = 2 }

# Changing the keys:
resource "persistent_counter" "example" {
  keys     = ["a", "b", "d", "c"]
}

# New result would be:
#   persistent_counter.example.values = { a = 0, b = 1, d = 3, c = 2 }

```

The provider is available from Terraform registry: [registry.terraform.io/providers/rosmo/persistent/latest](https://registry.terraform.io/providers/rosmo/persistent/latest)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.18

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```


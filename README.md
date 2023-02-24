# Persistent Counter provider for Terraform

```
resource "persistent_counter" "example" {
  keys     = ["a", "b", "c"]
}

persistent_counter.example.values = { a = 0, b = 1, c = 1}
```

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


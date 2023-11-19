# Persistent Counter provider for Terraform

This provider allows to assign unique, increasing integers to string keys
maintaining the same value per key while the key stays assigned.

## Example

```terraform
terraform {
  required_providers {
    persistent = {
      source  = "rosmo/persistent"
      version = ">=0.1.11"
    }
  }
}

variable input {
  type = set(string)
}

resource "persistent_counter" "with-reuse" {
  initial_value = 1
  keys          = var.input
  reuse         = true
}

resource "persistent_counter" "without-reuse" {
  initial_value = 1
  keys          = var.input
  reuse         = false
}

output counters {
  value = {
    "with-reuse" = persistent_counter.with-reuse.values,
    "without-reuse" = persistent_counter.without-reuse.values,
  }
}
```

Run the example providing an initial value as input results in the following output:

```shell
$ terraform apply -auto-approve -var 'input=["c","b","a"]'

counters = {
  "with-reuse" = tomap({
    "a" = 1
    "b" = 2
    "c" = 3
  })
  "without-reuse" = tomap({
    "a" = 1
    "b" = 2
    "c" = 3
  })
}
```

If values were just added, both, the version with and without `reuse` enabled
behave the same. Also note, that keys are always sorted ascending, before counter
values are assigned to them.

The difference the two versions in this example can be made clear when exchanging
an element (i.e. removing a value and adding a new one at the same time):

```shell
$ terraform apply -auto-approve -var 'input=["c","a","d"]'

counters = {
  "with-reuse" = tomap({
    "a" = 1
    "c" = 3
    "d" = 2
  })
  "without-reuse" = tomap({
    "a" = 1
    "c" = 3
    "d" = 4
  })
}
```

When `reuse` is set to `true`, the counter will re-assign values that are no
longer in use, while a value of `false` will always emit unique, ascending values.

## Usage

The provider is available from Terraform registry: [registry.terraform.io/providers/rosmo/persistent/latest](https://registry.terraform.io/providers/rosmo/persistent/latest).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://go.dev/doc/install) >= 1.21 (building from source only)

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

# Persistent Counter/Bucket provider for Terraform

This provider supports two resources:

- `persistent_counter`: assign unique, increasing integers to string keys maintaining the same value per key while the key stays assigned.
- `persistent_buckets`: fill buckets with keys according to capacity of each item

## Buckets example

```terraform
terraform {
  required_providers {
    persistent = {
      source  = "rosmo/persistent"
      version = ">=0.2.0"
    }
  }
}

resource "persistent_buckets" "example" {
  bucket_capacity = 100
  maximum_buckets = 2
  items = {
    item-1 = {
		  weight = 50
		  item   = "some string data here"
	  } 
    item-2 = {
		  weight = 25
		  item   = null
	  } 
    item-3 = {
		  weight = 10
	  } 
    item-4 = {
		  weight = 50
	  } 
    item-5 = {
		  weight = 10
	  } 
  }
}

# Result: 
# persistent_buckets.example.buckets = [
#   {  
#     item-1 = {
#  		  weight = 50
# 		  item   = "some string data here"
# 	  } 
#     item-2 = {
#   	  weight = 25
# 		  item   = null
#     } 
#     item-3 = {
# 	   	weight = 10
#	    } 
#     item-5 = {
#	  	  weight = 10
# 	  } 
#  },
#  {
#    item-4 = {
#		   weight = 50
#	   } 
#  } 
# ]
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

## Counter example

```terraform
terraform {
  required_providers {
    persistent = {
      source  = "rosmo/persistent"
      version = ">=0.2.0"
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

output "counters" {
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

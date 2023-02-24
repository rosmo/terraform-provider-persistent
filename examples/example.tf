terraform {
  required_providers {
    persistent-counter = {
      source  = "rosmo/persistent-counter"
      version = "0.1.0"
    }
  }
}

provider "persistent-counter" {
}

resource "persistent_counter_values" "example" {
  keys = ["a", "b", "d"]
}

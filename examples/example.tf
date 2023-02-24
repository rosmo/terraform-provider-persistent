terraform {
  required_providers {
    persistent = {
      source = "rosmo/persistent"
    }
  }
}

resource "persistent_counter" "example" {
  keys = ["a", "b", "d"]
}

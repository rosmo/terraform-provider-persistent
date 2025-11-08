resource "persistent_counter" "example" {
  initial_value = 5
  keys          = ["a", "b", "d"]
}
# persistent_counter.example.values = { a = 5, b = 6, c = 7 }

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPersistentCounterResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCounterResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("persistent_counter_values.test", "initial_value", "0"),
					resource.TestCheckResourceAttr("persistent_counter_values.test", "last_value", "2"),
					resource.TestCheckResourceAttr("persistent_counter_values.test", "values.a", "0"),
					resource.TestCheckResourceAttr("persistent_counter_values.test", "values.b", "1"),
					resource.TestCheckResourceAttr("persistent_counter_values.test", "values.c", "2"),
				),
			},
			{
				Config: testAccCounterAddResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("persistent_counter_values.test", "initial_value", "0"),
					resource.TestCheckResourceAttr("persistent_counter_values.test", "last_value", "3"),
					resource.TestCheckResourceAttr("persistent_counter_values.test", "values.a", "0"),
					resource.TestCheckResourceAttr("persistent_counter_values.test", "values.b", "1"),
					resource.TestCheckResourceAttr("persistent_counter_values.test", "values.c", "2"),
					resource.TestCheckResourceAttr("persistent_counter_values.test", "values.d", "3"),
				),
			},
			{
				Config: testAccCounterUpdateResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("persistent_counter_values.test", "initial_value", "0"),
					resource.TestCheckResourceAttr("persistent_counter_values.test", "last_value", "4"),
					resource.TestCheckResourceAttr("persistent_counter_values.test", "values.a", "0"),
					resource.TestCheckResourceAttr("persistent_counter_values.test", "values.b", "1"),
					resource.TestCheckResourceAttr("persistent_counter_values.test", "values.d", "3"),
					resource.TestCheckResourceAttr("persistent_counter_values.test", "values.e", "4"),
				),
			},
		},
	})
}

func testAccCounterResourceConfig() string {
	return `
resource "persistent_counter_values" "test" {
  provider      = persistent-counter
  initial_value = 0
  keys          = ["a", "b", "c"]
}
`
}

func testAccCounterAddResourceConfig() string {
	return `
resource "persistent_counter_values" "test" {
  provider      = persistent-counter
  initial_value = 0
  keys          = ["a", "b", "c", "d"]
}
`
}

func testAccCounterUpdateResourceConfig() string {
	return `
resource "persistent_counter_values" "test" {
  provider      = persistent-counter
  initial_value = 0
  keys          = ["a", "b", "d", "e"]
}
`
}

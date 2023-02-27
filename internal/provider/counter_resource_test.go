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
					resource.TestCheckResourceAttr("persistent_counter.test", "initial_value", "0"),
					resource.TestCheckResourceAttr("persistent_counter.test", "last_value", "2"),
					resource.TestCheckResourceAttr("persistent_counter.test", "values.a", "0"),
					resource.TestCheckResourceAttr("persistent_counter.test", "values.b", "1"),
					resource.TestCheckResourceAttr("persistent_counter.test", "values.c", "2"),
				),
			},
			{
				Config: testAccCounterAddResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("persistent_counter.test", "initial_value", "0"),
					resource.TestCheckResourceAttr("persistent_counter.test", "last_value", "3"),
					resource.TestCheckResourceAttr("persistent_counter.test", "values.a", "0"),
					resource.TestCheckResourceAttr("persistent_counter.test", "values.b", "1"),
					resource.TestCheckResourceAttr("persistent_counter.test", "values.c", "2"),
					resource.TestCheckResourceAttr("persistent_counter.test", "values.d", "3"),
				),
			},
			{
				Config: testAccCounterUpdateResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("persistent_counter.test", "initial_value", "0"),
					resource.TestCheckResourceAttr("persistent_counter.test", "last_value", "4"),
					resource.TestCheckResourceAttr("persistent_counter.test", "values.a", "0"),
					resource.TestCheckResourceAttr("persistent_counter.test", "values.b", "1"),
					resource.TestCheckResourceAttr("persistent_counter.test", "values.d", "3"),
					resource.TestCheckResourceAttr("persistent_counter.test", "values.e", "4"),
				),
			},
			{
				Config: testAccCounterInitialResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("persistent_counter.initial-value", "initial_value", "5"),
					resource.TestCheckResourceAttr("persistent_counter.initial-value", "last_value", "8"),
					resource.TestCheckResourceAttr("persistent_counter.initial-value", "values.a", "5"),
					resource.TestCheckResourceAttr("persistent_counter.initial-value", "values.b", "6"),
					resource.TestCheckResourceAttr("persistent_counter.initial-value", "values.c", "7"),
					resource.TestCheckResourceAttr("persistent_counter.initial-value", "values.d", "8"),
				),
			},
		},
	})
}

func testAccCounterResourceConfig() string {
	return `
resource "persistent_counter" "test" {
  initial_value = 0
  keys          = ["a", "b", "c"]
}
`
}

func testAccCounterAddResourceConfig() string {
	return `
resource "persistent_counter" "test" {
  initial_value = 0
  keys          = ["a", "b", "c", "d"]
}
`
}

func testAccCounterUpdateResourceConfig() string {
	return `
resource "persistent_counter" "test" {
  initial_value = 0
  keys          = ["a", "b", "d", "e"]
}
`
}

func testAccCounterInitialResourceConfig() string {
	return `
resource "persistent_counter" "initial-value" {
  initial_value = 5
  keys          = ["a", "b", "c", "d"]
}
`
}

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPersistentBucketsResource(t *testing.T) {
	errorRe, err := regexp.Compile("unable to find bucket capacity")
	if err != nil {
		panic(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBucketsResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("persistent_buckets.test", "bucket_capacity", "100"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "maximum_buckets", "1"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.#", "1"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.%", "3"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.weight", "50"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.item", "some string data here"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.weight", "25"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.item", ""),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-3.weight", "10"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-3.item", ""),
				),
			},
			{
				Config: testAccBucketsResourceAddConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("persistent_buckets.test", "bucket_capacity", "100"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "maximum_buckets", "1"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.#", "1"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.%", "4"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.weight", "50"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.item", "some string data here"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.weight", "25"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.item", ""),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-3.weight", "10"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-3.item", ""),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-4.weight", "10"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-4.item", ""),
				),
			},
			{
				Config: testAccBucketsResourceRemoveConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("persistent_buckets.test", "bucket_capacity", "100"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "maximum_buckets", "1"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.#", "1"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.%", "3"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.weight", "50"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.item", "some string data here"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-3.weight", "10"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-3.item", ""),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-4.weight", "10"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-4.item", ""),
				),
			},
			{
				Config:      testAccBucketsResourceTooMuchConfig(),
				ExpectError: errorRe,
			},
		},
	})
}

func TestAccPersistentBucketsGrowResource(t *testing.T) {
	errorRe, err := regexp.Compile("unable to find bucket capacity")
	if err != nil {
		panic(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBucketsResourceGrowConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("persistent_buckets.test", "bucket_capacity", "100"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "maximum_buckets", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.#", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.%", "3"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.weight", "50"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.item", "some string data here"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.weight", "25"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.item", ""),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-3.weight", "10"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-3.item", ""),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.%", "0"),
				),
			},
			{
				Config: testAccBucketsResourceGrowMoreConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("persistent_buckets.test", "bucket_capacity", "100"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "maximum_buckets", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.#", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.%", "3"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.weight", "50"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.item", "some string data here"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.weight", "25"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.item", ""),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-3.weight", "10"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-3.item", ""),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.%", "1"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.item-4.weight", "50"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.item-4.item", ""),
				),
			},
			{
				Config: testAccBucketsResourceGrowEvenMoreConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("persistent_buckets.test", "bucket_capacity", "100"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "maximum_buckets", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.#", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.%", "4"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.weight", "50"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.item", "some string data here"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.weight", "25"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.item", ""),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-3.weight", "10"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-3.item", ""),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-5.weight", "10"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-5.item", ""),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.%", "1"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.item-4.weight", "50"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.item-4.item", ""),
				),
			},
			{
				Config:      testAccBucketsResourceGrowTooMuchConfig(),
				ExpectError: errorRe,
			},
		},
	})
}

func TestAccPersistentBucketsTargetCapacityResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketsResourceTargetCapacityConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("persistent_buckets.test", "bucket_capacity", "100"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "target_capacity", "80"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "maximum_buckets", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.#", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.%", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.weight", "50"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.item", "some string data here"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.weight", "30"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.item", ""),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.%", "1"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.item-3.weight", "10"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.item-3.item", ""),
				),
			},
			{
				Config: testAccBucketsResourceTargetCapacityUpdateConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("persistent_buckets.test", "bucket_capacity", "100"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "target_capacity", "80"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "maximum_buckets", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.#", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.%", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.weight", "70"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.item", "some string data here"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.weight", "30"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.item", ""),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.%", "1"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.item-3.weight", "10"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.item-3.item", ""),
				),
			},
		},
	})
}

func TestAccPersistentBucketsMoveItemResource(t *testing.T) {
	errorRe, err := regexp.Compile("unable to find bucket capacity")
	if err != nil {
		panic(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketsResourceMoveItemConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("persistent_buckets.test", "bucket_capacity", "100"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "target_capacity", "80"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "move_items", "true"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "maximum_buckets", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.#", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.%", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.weight", "40"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-1.item", "some string data here"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.weight", "30"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.item", ""),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.%", "1"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.item-3.weight", "20"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.item-3.item", ""),
				),
			},
			{
				Config: testAccBucketsResourceMoveItemUpdateConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("persistent_buckets.test", "bucket_capacity", "100"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "target_capacity", "80"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "move_items", "true"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "maximum_buckets", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.#", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.%", "1"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.weight", "30"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.0.item-2.item", ""),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.%", "2"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.item-1.weight", "80"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.item-1.item", "some string data here"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.item-3.weight", "10"),
					resource.TestCheckResourceAttr("persistent_buckets.test", "buckets.1.item-3.item", ""),
				),
			},
			{
				Config:      testAccBucketsResourceMoveItemUpdateNoMoveConfig(),
				ExpectError: errorRe,
			},
		},
	})
}

func testAccBucketsResourceConfig() string {
	return `
resource "persistent_buckets" "test" {
  bucket_capacity = 100
  maximum_buckets = 1
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
  }
}
`
}

func testAccBucketsResourceAddConfig() string {
	return `
resource "persistent_buckets" "test" {
  bucket_capacity = 100
  maximum_buckets = 1
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
		weight = 10
	} 
  }
}
`
}

func testAccBucketsResourceRemoveConfig() string {
	return `
resource "persistent_buckets" "test" {
  bucket_capacity = 100
  maximum_buckets = 1
  items = {
    item-1 = {
		weight = 50
		item   = "some string data here"
	} 
    item-3 = {
		weight = 10
	} 
    item-4 = {
		weight = 10
	} 
  }
}
`
}

func testAccBucketsResourceTooMuchConfig() string {
	return `
resource "persistent_buckets" "test" {
  bucket_capacity = 100
  maximum_buckets = 1
  items = {
    item-1 = {
		weight = 50
		item   = "some string data here"
	} 
    item-3 = {
		weight = 10
	} 
    item-4 = {
		weight = 10
	} 
    item-5 = {
		weight = 50
	} 
  }
}
`
}

func testAccBucketsResourceGrowConfig() string {
	return `
resource "persistent_buckets" "test" {
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
  }
}
`
}

func testAccBucketsResourceGrowMoreConfig() string {
	return `
resource "persistent_buckets" "test" {
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
  }
}
`
}

func testAccBucketsResourceGrowEvenMoreConfig() string {
	return `
resource "persistent_buckets" "test" {
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
`
}
func testAccBucketsResourceGrowTooMuchConfig() string {
	return `
resource "persistent_buckets" "test" {
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
		weight = 125
	} 
  }
}
`
}

func testAccBucketsResourceTargetCapacityConfig() string {
	return `
resource "persistent_buckets" "test" {
  bucket_capacity = 100
  target_capacity = 80 
  maximum_buckets = 2
  items = {
    item-1 = {
		weight = 50
		item   = "some string data here"
	} 
    item-2 = {
		weight = 30
		item   = null
	} 
    item-3 = {
		weight = 10
	} 
  }
}
`
}

func testAccBucketsResourceTargetCapacityUpdateConfig() string {
	return `
resource "persistent_buckets" "test" {
  bucket_capacity = 100
  target_capacity = 80 
  maximum_buckets = 2
  items = {
    item-1 = {
		weight = 70
		item   = "some string data here"
	} 
    item-2 = {
		weight = 30
		item   = null
	} 
    item-3 = {
		weight = 10
	} 
  }
}
`
}

func testAccBucketsResourceMoveItemConfig() string {
	return `
resource "persistent_buckets" "test" {
  bucket_capacity = 100
  target_capacity = 80 
  maximum_buckets = 2
  items = {
    item-1 = {
		weight = 40
		item   = "some string data here"
	} 
    item-2 = {
		weight = 30
		item   = null
	} 
    item-3 = {
		weight = 20
	} 
  }
}
`
}

func testAccBucketsResourceMoveItemUpdateConfig() string {
	return `
resource "persistent_buckets" "test" {
  bucket_capacity = 100
  target_capacity = 80 
  maximum_buckets = 2
  move_items      = true
  items = {
    item-1 = {
		weight = 80
		item   = "some string data here"
	} 
    item-2 = {
		weight = 30
		item   = null
	} 
    item-3 = {
		weight = 10
	} 
  }
}
`
}

func testAccBucketsResourceMoveItemUpdateNoMoveConfig() string {
	return `
resource "persistent_buckets" "test" {
  bucket_capacity = 100
  target_capacity = 80 
  maximum_buckets = 2
  move_items      = false
  items = {
    item-1 = {
		weight = 80
		item   = "some string data here"
	} 
    item-2 = {
		weight = 30
		item   = null
	} 
    item-3 = {
		weight = 30
	} 
  }
}
`
}

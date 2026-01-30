package storage_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vpsie/terraform-provider-vpsie/internal/acctest"
)

func TestAccStorageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create and read testing
			{
				Config: testAccVpsieStorageConfig_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vpsie_storage.test", "dc_identifier", "eu-west-1"),
					resource.TestCheckResourceAttr("vpsie_storage.test", "size", "10"),
					resource.TestCheckResourceAttr("vpsie_storage.test", "description", "test description"),
					resource.TestCheckResourceAttr("vpsie_storage.test", "name", "tf-acc-test-storage"),
					resource.TestCheckResourceAttr("vpsie_storage.test", "disk_format", "EXT4"),
					resource.TestCheckResourceAttr("vpsie_storage.test", "storage_type", "standard"),
					resource.TestCheckResourceAttrSet("vpsie_storage.test", "identifier"),
				),
			},
			{
				ResourceName:      "vpsie_storage.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// update size
			{
				Config: testAccVpsieStorageConfig_update("tf-acc-test-storage", "20"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vpsie_storage.test", "dc_identifier", "eu-west-1"),
					resource.TestCheckResourceAttr("vpsie_storage.test", "size", "20"),
					resource.TestCheckResourceAttr("vpsie_storage.test", "description", "test description"),
					resource.TestCheckResourceAttr("vpsie_storage.test", "name", "tf-acc-test-storage"),
				),
			},
		},
	})
}

func TestAccStorageResource_Resize(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create and read testing
			{
				Config: testAccVpsieStorageConfig_update("tf-acc-test-storage-resize", "20"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vpsie_storage.test", "dc_identifier", "eu-west-1"),
					resource.TestCheckResourceAttr("vpsie_storage.test", "description", "test description"),
					resource.TestCheckResourceAttr("vpsie_storage.test", "name", "tf-acc-test-storage-resize"),
					resource.TestCheckResourceAttr("vpsie_storage.test", "size", "20"),
				),
			},
			{
				Config: testAccVpsieStorageConfig_update("tf-acc-test-storage-resize", "50"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vpsie_storage.test", "dc_identifier", "eu-west-1"),
					resource.TestCheckResourceAttr("vpsie_storage.test", "description", "test description"),
					resource.TestCheckResourceAttr("vpsie_storage.test", "name", "tf-acc-test-storage-resize"),
					resource.TestCheckResourceAttr("vpsie_storage.test", "size", "50"),
				),
			},
		},
	})
}

const testAccVpsieStorageConfig_basic = `
resource "vpsie_storage" "test" {
  dc_identifier = "eu-west-1"
  name          = "tf-acc-test-storage"
  size          = 10
  storage_type  = "standard"
  disk_format   = "EXT4"
  description   = "test description"
}
`

func testAccVpsieStorageConfig_update(name, size string) string {
	return fmt.Sprintf(`
resource "vpsie_storage" "test" {
  dc_identifier = "eu-west-1"
  name          = %q
  size          = %s
  storage_type  = "standard"
  disk_format   = "EXT4"
  description   = "test description"
}
`, name, size)
}

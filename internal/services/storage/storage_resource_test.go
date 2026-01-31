package storage_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vpsie/govpsie"
	"github.com/vpsie/terraform-provider-vpsie/internal/acctest"
)

func TestAccStorageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageResourceDestroy,
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
		CheckDestroy:             testAccCheckStorageResourceDestroy,
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

func testAccCheckStorageResourceDestroy(s *terraform.State) error {
	client := govpsie.NewClient(nil)
	client.SetRequestHeaders(map[string]string{
		"Vpsie-Auth": os.Getenv("VPSIE_ACCESS_TOKEN"),
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vpsie_storage" {
			continue
		}

		volumes, err := client.Storage.ListAll(context.Background(), &govpsie.ListOptions{})
		if err != nil {
			return fmt.Errorf("error checking storage destroy: %s", err)
		}

		for _, volume := range volumes {
			if volume.Identifier == rs.Primary.Attributes["identifier"] {
				return fmt.Errorf("storage volume %s still exists", rs.Primary.Attributes["identifier"])
			}
		}
	}

	return nil
}

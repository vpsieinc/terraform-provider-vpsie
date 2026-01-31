package project_test

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

func TestAccProjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProjectResourceDestroy,
		Steps: []resource.TestStep{
			// Create and read
			{
				Config: testAccProjectConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vpsie_project.test", "name", "tf-acc-test-project"),
					resource.TestCheckResourceAttr("vpsie_project.test", "description", "acceptance test project"),
					resource.TestCheckResourceAttrSet("vpsie_project.test", "identifier"),
					resource.TestCheckResourceAttrSet("vpsie_project.test", "id"),
					resource.TestCheckResourceAttrSet("vpsie_project.test", "created_on"),
				),
			},
			// Import
			{
				ResourceName:      "vpsie_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccProjectConfig = `
resource "vpsie_project" "test" {
  name        = "tf-acc-test-project"
  description = "acceptance test project"
}
`

func testAccCheckProjectResourceDestroy(s *terraform.State) error {
	client := govpsie.NewClient(nil)
	client.SetRequestHeaders(map[string]string{
		"Vpsie-Auth": os.Getenv("VPSIE_ACCESS_TOKEN"),
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vpsie_project" {
			continue
		}

		projects, err := client.Project.List(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("error checking project destroy: %s", err)
		}

		for _, project := range projects {
			if project.Identifier == rs.Primary.Attributes["identifier"] {
				return fmt.Errorf("project %s still exists", rs.Primary.Attributes["identifier"])
			}
		}
	}

	return nil
}

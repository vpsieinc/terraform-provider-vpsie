package project_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vpsie/terraform-provider-vpsie/internal/acctest"
)

func TestAccProjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
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

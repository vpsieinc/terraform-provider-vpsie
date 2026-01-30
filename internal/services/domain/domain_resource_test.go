package domain_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vpsie/terraform-provider-vpsie/internal/acctest"
)

func TestAccDomainResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and read
			{
				Config: testAccDomainConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vpsie_domain.test", "domain_name", "tf-acc-test.example.com"),
					resource.TestCheckResourceAttrSet("vpsie_domain.test", "identifier"),
					resource.TestCheckResourceAttrSet("vpsie_domain.test", "created_on"),
					resource.TestCheckResourceAttrSet("vpsie_domain.test", "project_identifier"),
				),
			},
			// Import
			{
				ResourceName:      "vpsie_domain.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDomainConfig = `
resource "vpsie_project" "test" {
  name        = "tf-acc-test-domain-project"
  description = "project for domain acceptance test"
}

resource "vpsie_domain" "test" {
  domain_name        = "tf-acc-test.example.com"
  project_identifier = vpsie_project.test.identifier
}
`

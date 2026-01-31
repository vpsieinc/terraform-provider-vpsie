package domain_test

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

func TestAccDomainResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDomainResourceDestroy,
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

func testAccCheckDomainResourceDestroy(s *terraform.State) error {
	client := govpsie.NewClient(nil)
	client.SetRequestHeaders(map[string]string{
		"Vpsie-Auth": os.Getenv("VPSIE_ACCESS_TOKEN"),
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vpsie_domain" {
			continue
		}

		domains, err := client.Domain.ListDomains(context.Background(), &govpsie.ListOptions{})
		if err != nil {
			return fmt.Errorf("error checking domain destroy: %s", err)
		}

		for _, domain := range domains {
			if domain.DomainName == rs.Primary.Attributes["domain_name"] {
				return fmt.Errorf("domain %s still exists", rs.Primary.Attributes["domain_name"])
			}
		}
	}

	return nil
}

package accesstoken_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vpsie/terraform-provider-vpsie/internal/acctest"
)

func TestAccAccessTokenResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and read
			{
				Config: testAccAccessTokenConfig("tf-acc-test-token", "2030-01-01"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vpsie_access_token.test", "name", "tf-acc-test-token"),
					resource.TestCheckResourceAttr("vpsie_access_token.test", "expiration_date", "2030-01-01"),
					resource.TestCheckResourceAttrSet("vpsie_access_token.test", "identifier"),
					resource.TestCheckResourceAttrSet("vpsie_access_token.test", "created_on"),
				),
			},
			// Update name and expiration
			{
				Config: testAccAccessTokenConfig("tf-acc-test-token-updated", "2031-01-01"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vpsie_access_token.test", "name", "tf-acc-test-token-updated"),
					resource.TestCheckResourceAttr("vpsie_access_token.test", "expiration_date", "2031-01-01"),
				),
			},
		},
	})
}

func testAccAccessTokenConfig(name, expiration string) string {
	return fmt.Sprintf(`
resource "vpsie_access_token" "test" {
  name            = %q
  access_token    = "tf-acc-test-token-value"
  expiration_date = %q
}
`, name, expiration)
}

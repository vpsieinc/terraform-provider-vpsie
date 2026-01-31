package accesstoken_test

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

func TestAccAccessTokenResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAccessTokenResourceDestroy,
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
  access_token    = "tf-acc-test-token-value-that-is-long-enough-to-meet-the-64-char-minimum"
  expiration_date = %q
}
`, name, expiration)
}

func testAccCheckAccessTokenResourceDestroy(s *terraform.State) error {
	client := govpsie.NewClient(nil)
	client.SetRequestHeaders(map[string]string{
		"Vpsie-Auth": os.Getenv("VPSIE_ACCESS_TOKEN"),
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vpsie_access_token" {
			continue
		}

		tokens, err := client.AccessToken.List(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("error checking access token destroy: %s", err)
		}

		for _, token := range tokens {
			if token.AccessTokenIdentifier == rs.Primary.Attributes["identifier"] {
				return fmt.Errorf("access token %s still exists", rs.Primary.Attributes["identifier"])
			}
		}
	}

	return nil
}

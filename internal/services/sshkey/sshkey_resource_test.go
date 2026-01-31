package sshkey_test

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

func TestAccSshkeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSshkeyResourceDestroy,
		Steps: []resource.TestStep{
			// Create and read
			{
				Config: testAccSshkeyConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vpsie_sshkey.test", "name", "tf-acc-test-sshkey"),
					resource.TestCheckResourceAttrSet("vpsie_sshkey.test", "identifier"),
					resource.TestCheckResourceAttrSet("vpsie_sshkey.test", "id"),
					resource.TestCheckResourceAttrSet("vpsie_sshkey.test", "created_on"),
				),
			},
			// Import
			{
				ResourceName:            "vpsie_sshkey.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key"},
			},
		},
	})
}

const testAccSshkeyConfig = `
resource "vpsie_sshkey" "test" {
  name        = "tf-acc-test-sshkey"
  private_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGtBFkE3xMDgPSGMgdCLvPqDC0yMU7gXgEBqifN6sXlu tf-acc-test"
}
`

func testAccCheckSshkeyResourceDestroy(s *terraform.State) error {
	client := govpsie.NewClient(nil)
	client.SetRequestHeaders(map[string]string{
		"Vpsie-Auth": os.Getenv("VPSIE_ACCESS_TOKEN"),
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vpsie_sshkey" {
			continue
		}

		sshkeys, err := client.SShKey.List(context.Background())
		if err != nil {
			return fmt.Errorf("error checking sshkey destroy: %s", err)
		}

		for _, sshkey := range sshkeys {
			if sshkey.Name == rs.Primary.Attributes["identifier"] {
				return fmt.Errorf("sshkey %s still exists", rs.Primary.Attributes["identifier"])
			}
		}
	}

	return nil
}

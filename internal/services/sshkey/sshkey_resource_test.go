package sshkey_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vpsie/terraform-provider-vpsie/internal/acctest"
)

func TestAccSshkeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
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

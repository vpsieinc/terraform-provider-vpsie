package script_test

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

func TestAccScriptResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckScriptResourceDestroy,
		Steps: []resource.TestStep{
			// Create and read
			{
				Config: testAccScriptConfig("tf-acc-test-script", "#!/bin/bash\necho hello"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vpsie_script.test", "script_name", "tf-acc-test-script"),
					resource.TestCheckResourceAttr("vpsie_script.test", "type", "bash"),
					resource.TestCheckResourceAttrSet("vpsie_script.test", "identifier"),
					resource.TestCheckResourceAttrSet("vpsie_script.test", "id"),
					resource.TestCheckResourceAttrSet("vpsie_script.test", "created_on"),
				),
			},
			// Import
			{
				ResourceName:      "vpsie_script.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update
			{
				Config: testAccScriptConfig("tf-acc-test-script-updated", "#!/bin/bash\necho updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("vpsie_script.test", "script_name", "tf-acc-test-script-updated"),
				),
			},
		},
	})
}

func testAccScriptConfig(name, script string) string {
	return fmt.Sprintf(`
resource "vpsie_script" "test" {
  script_name = %q
  script      = %q
  type        = "bash"
}
`, name, script)
}

func testAccCheckScriptResourceDestroy(s *terraform.State) error {
	client := govpsie.NewClient(nil)
	client.SetRequestHeaders(map[string]string{
		"Vpsie-Auth": os.Getenv("VPSIE_ACCESS_TOKEN"),
	})

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vpsie_script" {
			continue
		}

		scripts, err := client.Scripts.GetScripts(context.Background())
		if err != nil {
			return fmt.Errorf("error checking script destroy: %s", err)
		}

		for _, script := range scripts {
			if script.Identifier == rs.Primary.Attributes["identifier"] {
				return fmt.Errorf("script %s still exists", rs.Primary.Attributes["identifier"])
			}
		}
	}

	return nil
}

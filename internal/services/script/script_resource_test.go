package script_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vpsie/terraform-provider-vpsie/internal/acctest"
)

func TestAccScriptResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
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

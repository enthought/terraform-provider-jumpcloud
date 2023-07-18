package jumpcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataApplication(t *testing.T) {
	const cName = "OpenID Connect"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApplication(cName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.jumpcloud_application.test_application", "name", cName),
				),
			},
		},
	})
}

func testAccDataApplication(name string) string {
	return fmt.Sprintf(`
		data "jumpcloud_application" "test_application" {
			name = "%s"
		}`, name,
	)
}

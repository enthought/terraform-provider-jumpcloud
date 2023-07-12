package jumpcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSystemGroup(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemGroup(rName),
				Check:  resource.TestCheckResourceAttr("jumpcloud_system_group.test_system_group", "name", rName),
			},
		},
	})
}

func testAccSystemGroup(name string) string {
	return fmt.Sprintf(`
		resource "jumpcloud_system_group" "test_system_group" {
    		name = "%s"
		}`, name,
	)
}

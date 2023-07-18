package jumpcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserGroup(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	posixName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	gid := acctest.RandIntRange(1, 1000)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroup(rName, gid, posixName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_user_group.test_group", "name", rName),
				),
			},
		},
	})
}

func testAccUserGroup(name string, gid int, posixName string) string {
	return fmt.Sprintf(`
		resource "jumpcloud_user_group" "test_group" {
    		name = "%s"
		}`, name,
	)
}

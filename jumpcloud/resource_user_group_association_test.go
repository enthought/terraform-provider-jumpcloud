package jumpcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserGroupAssociation(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupAssociation(rName),
				Check: resource.TestCheckResourceAttrSet("jumpcloud_user_group_association.test_user_group_association_"+rName,
					"appid"),
			},
		},
	})
}

func testAccUserGroupAssociation(name string) string {
	return fmt.Sprintf(`
		data "jumpcloud_application" "test_application_%s" {
			name = "OpenID Connect"
		}

		resource "jumpcloud_user_group" "test_user_group_%s" {
			name = "testgroup_%s"
		}

		resource "jumpcloud_user_group_association" "test_user_group_association_%s" {
			groupid = jumpcloud_user_group.test_user_group_%s.id
			appid = data.jumpcloud_application.test_application_%s.id
		}
  `, name, name, name, name, name, name)
}

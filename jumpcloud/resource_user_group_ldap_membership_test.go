package jumpcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserGroupLdapMembership(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupLdapMembership(rName),
				Check: resource.TestCheckResourceAttrSet("jumpcloud_user_group_ldap_membership.test_user_group_ldap_membership_"+rName,
					"ldap_id"),
			},
		},
	})
}

func testAccUserGroupLdapMembership(name string) string {
	return fmt.Sprintf(`
		resource "jumpcloud_user_group" "test_user_group_%s" {
			name = "testgroup_%s"
		}

		data "jumpcloud_ldap_server" "test_ldap_server_%s" {
			name="jumpcloud" 
		}

		resource "jumpcloud_user_group_ldap_membership" "test_user_group_ldap_membership_%s" {
			usergroup_id = jumpcloud_user_group.test_user_group_%s.id
			ldap_id = data.jumpcloud_ldap_server.test_ldap_server_%s.ldap_id
		}
	`, name, name, name, name, name, name)
}

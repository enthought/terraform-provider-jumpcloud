package jumpcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataLdapServer(t *testing.T) {
	const cName = "jumpcloud"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDataLdapServer(cName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.jumpcloud_ldap_server.test_ldap_server", "name", cName),
				),
			},
		},
	})
}

func testAccDataLdapServer(name string) string {
	return fmt.Sprintf(`
		data "jumpcloud_ldap_server" "test_ldap_server" {
			name = "%s"
		}`, name,
	)
}

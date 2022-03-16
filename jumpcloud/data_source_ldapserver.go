package jumpcloud

import (
	"context"
	"fmt"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataResourceLdapServer() *schema.Resource {
	return &schema.Resource{
		Read: dataResourceLdapServerRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_lockout_action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_password_expiration_action": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataResourceLdapServerRead(d *schema.ResourceData, m interface{}) error {

	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	var ldap_filter []string

	if len(d.Get("ldap_id").(string)) > 0 {
		ldap_filter = append(ldap_filter, "id:eq:"+d.Get("ldap_id").(string))
	} else if len(d.Get("name").(string)) > 0 {
		ldap_filter = append(ldap_filter, "name:eq:"+d.Get("name").(string))
	} else {
		return fmt.Errorf("ldap_id or name must be set for jumpcloud_ldap_server")
	}

	payload := map[string]interface{}{
		"filter": ldap_filter,
	}

	req := map[string]interface{}{
		"body": payload,
	}

	res, _, err := client.LDAPServersApi.LdapserversList(context.TODO(), "", headerAccept, req)
	if err != nil {
		return err
	}

	if err := d.Set("name", res[0].Name); err != nil {
		return err
	}

	if err := d.Set("ldap_id", res[0].Id); err != nil {
		return err
	}

	d.SetId(d.Get("ldap_id").(string))

	if err := d.Set("user_lockout_action", res[0].UserLockoutAction); err != nil {
		return err
	}

	if err := d.Set("user_password_expiration_action", res[0].UserPasswordExpirationAction); err != nil {
		return err
	}

	return nil
}

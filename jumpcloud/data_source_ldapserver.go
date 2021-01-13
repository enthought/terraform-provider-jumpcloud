package jumpcloud

import (
	"context"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataResourceLdapServer() *schema.Resource {
	return &schema.Resource{
		Read: dataResourceLdapServerRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Required: true,
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

	filter_by_name := [...]string{"name:eq:" + d.Get("name").(string)}

	payload := map[string]interface{}{
		"filter": filter_by_name,
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

	if err := d.Set("id", res[0].Id); err != nil {
		return err
	}

	if err := d.Set("user_lockout_action", res[0].UserLockoutAction); err != nil {
		return err
	}

	if err := d.Set("user_password_expiration_action", res[0].UserPasswordExpirationAction); err != nil {
		return err
	}

	return nil
}

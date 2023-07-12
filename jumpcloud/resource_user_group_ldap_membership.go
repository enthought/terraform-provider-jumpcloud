package jumpcloud

import (
	"context"
	"strings"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserGroupLdapMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserGroupLdapMembershipCreate,
		Read:   resourceUserGroupLdapMembershipRead,
		Delete: resourceUserGroupLdapMembershipDelete,
		Schema: map[string]*schema.Schema{
			"ldap_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"usergroup_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceUserGroupLdapMembershipCreate(d *schema.ResourceData, m interface{}) error {

	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	ugId := d.Get("usergroup_id").(string)
	ldapId := d.Get("ldap_id").(string)

	payload := jcapiv2.UserGroupGraphManagementReq{
		Id:    ldapId,
		Op:    "add",
		Type_: "ldap_server",
	}

	req := map[string]interface{}{
		"body": payload,
	}

	_, err := client.UserGroupsApi.GraphUserGroupAssociationsPost(context.TODO(),
		ugId, "", headerAccept, req)

	if err != nil {
		return err
	}

	d.SetId(ugId + ":" + ldapId)
	return resourceUserGroupLdapMembershipRead(d, m)

}

func resourceUserGroupLdapMembershipRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	ids := strings.Split(d.Id(), ":")

	ugId := ids[0]
	ldapId := [...]string{ids[1]}

	req := map[string]interface{}{
		"targets": ldapId,
	}

	_, _, err := client.UserGroupsApi.GraphUserGroupTraverseLdapServer(context.TODO(),
		ugId, "", "", req)

	if err != nil {
		return err
	}

	if err := d.Set("usergroup_id", ugId); err != nil {
		return err
	}
	if err := d.Set("ldap_id", ldapId[0]); err != nil {
		return err
	}
	return nil
}

func resourceUserGroupLdapMembershipDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	ids := strings.Split(d.Id(), ":")

	ugId := ids[0]
	ldapId := ids[1]

	payload := jcapiv2.UserGroupGraphManagementReq{
		Id:    ldapId,
		Op:    "remove",
		Type_: "ldap_server",
	}

	req := map[string]interface{}{
		"body": payload,
	}

	_, err := client.UserGroupsApi.GraphUserGroupAssociationsPost(context.TODO(),
		ugId, "", headerAccept, req)

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

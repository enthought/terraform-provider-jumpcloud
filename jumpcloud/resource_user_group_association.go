package jumpcloud

import (
	"context"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserGroupAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserGroupAssociationCreate,
		Read:   resourceUserGroupAssociationRead,
		Update: nil,
		Delete: resourceUserGroupAssociationDelete,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Description: "The ID of the `resource_user_group` resource.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"object_id": {
				Description: "The ID of the object to associate to the group.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func modifyUserGroupAssociation(client *jcapiv2.APIClient,
	d *schema.ResourceData, action string) error {

	payload := jcapiv2.UserGroupGraphManagementReq{
		Op:    action,
		Type_: "application",
		Id:    d.Get("object_id").(string),
	}

	req := map[string]interface{}{
		"body": payload,
	}

	_, err := client.UserGroupAssociationsApi.GraphUserGroupAssociationsPost(
		context.TODO(), d.Get("group_id").(string), "", "", req)

	return err
}

func resourceUserGroupAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	err := modifyUserGroupAssociation(client, d, "add")
	if err != nil {
		return err
	}
	return resourceUserGroupAssociationRead(d, meta)
}

func resourceUserGroupAssociationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	optionals := map[string]interface{}{
		"groupId": d.Get("group_id").(string),
		"limit":   int32(100),
	}

	graphconnect, _, err := client.UserGroupAssociationsApi.GraphUserGroupAssociationsList(
		context.TODO(), d.Get("group_id").(string), "", "", []string{"application"}, optionals)
	if err != nil {
		return err
	}

	for _, v := range graphconnect {
		if v.To.Id == d.Get("object_id") {
			resourceId := d.Get("group_id").(string) + "/" + d.Get("object_id").(string)
			d.SetId(resourceId)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceUserGroupAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)
	return modifyUserGroupAssociation(client, d, "remove")
}

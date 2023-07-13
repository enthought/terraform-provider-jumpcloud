package jumpcloud

import (
	"context"
	"strings"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserGroupAssociation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserGroupAssociationCreate,
		ReadContext:   resourceUserGroupAssociationRead,
		UpdateContext: nil,
		DeleteContext: resourceUserGroupAssociationDelete,
		Schema: map[string]*schema.Schema{
			"groupid": {
				Description: "The ID of the `resource_user_group` resource.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"appid": {
				Description: "The ID of the application to associate to the group.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: userGroupAssociationImporter,
		},
	}
}

func userGroupAssociationImporter(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), "/")
	d.Set("groupid", s[0])
	d.Set("appid", s[1])
	return []*schema.ResourceData{d}, nil
}

func modifyUserGroupAssociation(ctx context.Context, client *jcapiv2.APIClient, d *schema.ResourceData, action string) error {
	payload := jcapiv2.UserGroupGraphManagementReq{
		Op:    action,
		Type_: "application",
		Id:    d.Get("appid").(string),
	}

	req := map[string]interface{}{
		"body": payload,
	}

	_, err := client.UserGroupAssociationsApi.GraphUserGroupAssociationsPost(ctx, d.Get("groupid").(string), "", "", req)

	return err
}

func resourceUserGroupAssociationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	err := modifyUserGroupAssociation(ctx, client, d, "add")
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceUserGroupAssociationRead(ctx, d, meta)
}

func resourceUserGroupAssociationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	optionals := map[string]interface{}{
		"groupId": d.Get("groupid").(string),
		"limit":   int32(100),
	}

	graphconnect, _, err := client.UserGroupAssociationsApi.GraphUserGroupAssociationsList(
		ctx, d.Get("groupid").(string), "", "", []string{"application"}, optionals)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range graphconnect {
		if v.To.Id == d.Get("appid") {
			resourceId := d.Get("groupid").(string) + "/" + d.Get("appid").(string)
			d.SetId(resourceId)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceUserGroupAssociationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	err := modifyUserGroupAssociation(ctx, client, d, "remove")

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

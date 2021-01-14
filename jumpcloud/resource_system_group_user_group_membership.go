package jumpcloud

import (
	"context"
	"fmt"
	"strings"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSystemGroupUserGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceSystemGroupUserGroupMembershipCreate,
		Read:   resourceSystemGroupUserGroupMembershipRead,
		Delete: resourceSystemGroupUserGroupMembershipDelete,
		Schema: map[string]*schema.Schema{
			"systems_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"users_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
				ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceSystemGroupUserGroupMembershipCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	body := jcapiv2.SystemGroupGraphManagementReq{
		Id:    d.Get("users_group_id").(string),
		Op:    "add",
		Type_: "user_group",
	}

	id := d.Get("systems_group_id").(string)

	optional := map[string]interface{}{
		"groupId": id,
		"body":    body,
	}
	res, err := client.GraphApi.GraphSystemGroupAssociationsPost(context.TODO(),
		id, "", headerAccept, optional)
	if err != nil {
		// TODO: sort out error essentials
		return fmt.Errorf("error creating system group association %+v: = %+v",
			err, res)
	}

	d.SetId(d.Get("systems_group_id").(string) + ":" + d.Get("users_group_id").(string))
	d.Set("active", true)
	return nil
}

func resourceSystemGroupUserGroupMembershipRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	split_id := strings.Split(d.Id(), ":")

	system_group_id := split_id[0]
	user_group_id := split_id[1]

	var targets []string

	targets = append(targets, "user_group")

	associations, res, err := client.GraphApi.GraphSystemAssociationsList(context.TODO(),
		system_group_id, "", headerAccept, targets, nil)
	if err != nil {
		// TODO: sort out error essentials
		return fmt.Errorf("error reading system group associations %s: %s - response = %+v",
			system_group_id, err, res)
	}

	d.Set("active", false)
	for _, association := range associations {
		if association.To.Id == user_group_id {
			d.Set("active", true)
		}
	}

	return nil
}

func resourceSystemGroupUserGroupMembershipDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	body := jcapiv2.SystemGroupGraphManagementReq{
		Id:    d.Get("users_group_id").(string),
		Op:    "remove",
		Type_: "user_group",
	}

	var id string
	id = d.Get("systems_group_id").(string)

	req := map[string]interface{}{
		"body": body,
	}
	res, err := client.GraphApi.GraphSystemGroupAssociationsPost(context.TODO(),
		id, "", headerAccept, req)
	if err != nil {
		// TODO: sort out error essentials
		return fmt.Errorf("error deleting system group %s: response = %+v",
			err, res)
	}

	d.SetId("")
	return nil
}

package jumpcloud

import (
	"context"
	"fmt"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroupsSystem() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupsSystemCreate,
		ReadContext:   resourceGroupsSystemRead,
		UpdateContext: resourceGroupsSystemUpdate,
		DeleteContext: resourceGroupsSystemDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceGroupsSystemCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	body := jcapiv2.SystemGroupData{Name: d.Get("name").(string)}

	req := map[string]interface{}{
		"body": body,
	}
	group, res, err := client.SystemGroupsApi.GroupsSystemPost(ctx, "", headerAccept, req)
	if err != nil {
		// TODO: sort out error essentials
		return diag.Errorf("error creating system group %s: %s - response = %+v",
			(req["body"].(jcapiv2.SystemGroupData)).Name, err, res)
	}

	d.SetId(group.Name)
	d.Set("name", group.Name)
	d.Set("jc_id", group.Id)
	return resourceGroupsSystemRead(ctx, d, m)
}

// Helper to look up a system group by name
func resourceGroupsSystemList_match(ctx context.Context, d *schema.ResourceData, m interface{}) (jcapiv2.SystemGroup, error) {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	var filter []string

	filter = append(filter, "name:eq:"+d.Id())

	optional := map[string]interface{}{
		"filter": filter,
	}

	result, _, err := client.SystemGroupsApi.GroupsSystemList(ctx, "", headerAccept, optional)
	if err == nil {
		if len(result) < 1 {
			return jcapiv2.SystemGroup{}, fmt.Errorf("system group \"%s\" not found", d.Id())
		} else {
			return result[0], nil
		}
	} else {
		return jcapiv2.SystemGroup{}, err
	}
}

func resourceGroupsSystemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	var id string = d.Get("jc_id").(string)

	if id == "" {
		id_lookup, err := resourceGroupsSystemList_match(ctx, d, m)
		if err != nil {
			return diag.Errorf("unable to locate ID for group %s, %+v", d.Get("name"), err)
		}
		id = id_lookup.Id
		d.SetId(id_lookup.Name)
		d.Set("name", id_lookup.Name)
		d.Set("jc_id", id_lookup.Id)
	}

	group, res, err := client.SystemGroupsApi.GroupsSystemGet(ctx, id, "", headerAccept, nil)
	if err != nil {
		// TODO: sort out error essentials
		return diag.Errorf("error reading system group ID %s: %s - response = %+v", d.Id(), err, res)
	}

	d.SetId(group.Name)
	d.Set("name", group.Name)
	d.Set("jc_id", group.Id)
	return nil
}

func resourceGroupsSystemUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	var id string = d.Get("jc_id").(string)

	body := jcapiv2.SystemGroupData{Name: d.Get("name").(string)}

	req := map[string]interface{}{
		"body": body,
	}

	group, res, err := client.SystemGroupsApi.GroupsSystemPut(ctx, id, "", headerAccept, req)
	if err != nil {
		// TODO: sort out error essentials
		return diag.Errorf("error updating system group %s: %s - response = %+v", d.Get("name"), err, res)
	}

	d.SetId(group.Name)
	d.Set("name", group.Name)
	d.Set("jc_id", group.Id)
	return resourceGroupsSystemRead(ctx, d, m)
}

func resourceGroupsSystemDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	var id string = d.Get("jc_id").(string)

	res, err := client.SystemGroupsApi.GroupsSystemDelete(ctx, id, "", headerAccept, nil)
	if err != nil {
		// TODO: sort out error essentials
		return diag.Errorf("error deleting system group:%s; response = %+v", err, res)
	}
	d.SetId("")
	return nil
}

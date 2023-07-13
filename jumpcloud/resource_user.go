package jumpcloud

import (
	"context"

	jcapiv1 "github.com/TheJumpCloud/jcapi-go/v1"
	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"firstname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lastname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_mfa": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// Currently, only the options necessary for our use case are implemented
			// JumpCloud offers a lot more
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	payload := jcapiv1.Systemuserputpost{
		Username:                    d.Get("username").(string),
		Email:                       d.Get("email").(string),
		Firstname:                   d.Get("firstname").(string),
		Lastname:                    d.Get("lastname").(string),
		Password:                    d.Get("password").(string),
		EnableUserPortalMultifactor: d.Get("enable_mfa").(bool),
	}
	req := map[string]interface{}{
		"body": payload,
	}
	returnstruc, _, err := client.SystemusersApi.SystemusersPost(ctx, "", "", req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(returnstruc.Id)
	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	res, _, err := client.SystemusersApi.SystemusersGet(ctx, d.Id(), "", "", nil)

	// If the object does not exist in our infrastructure, we unset the ID
	// Unfortunately, the http request returns 200 even if the resource does not exist
	if err != nil {
		if err.Error() == "EOF" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId(res.Id)

	if err := d.Set("username", res.Username); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email", res.Email); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("firstname", res.Firstname); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("lastname", res.Lastname); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enable_mfa", res.EnableUserPortalMultifactor); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	// The code from the create function is almost identical, but the structure is different :
	// jcapiv1.Systemuserput != jcapiv1.Systemuserputpost
	payload := jcapiv1.Systemuserput{
		Username:                    d.Get("username").(string),
		Email:                       d.Get("email").(string),
		Firstname:                   d.Get("firstname").(string),
		Lastname:                    d.Get("lastname").(string),
		Password:                    d.Get("password").(string),
		EnableUserPortalMultifactor: d.Get("enable_mfa").(bool),
	}

	req := map[string]interface{}{
		"body": payload,
	}
	_, _, err := client.SystemusersApi.SystemusersPut(ctx, d.Id(), "", "", req)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	res, _, err := client.SystemusersApi.SystemusersDelete(ctx, d.Id(), "", headerAccept, nil)
	if err != nil {
		// TODO: sort out error essentials
		return diag.Errorf("error deleting user:%s; response = %+v", err, res)
	}
	d.SetId("")
	return nil
}

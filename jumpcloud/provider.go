package jumpcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider instantiates a terraform provider for Jumpcloud
// This includes all operations on all supported resources and
// global Jumpcloud parameters
func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUMPCLOUD_API_KEY", nil),
				Description: descriptions["api_key"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"jumpcloud_user":                  resourceUser(),
			"jumpcloud_user_group":            resourceUserGroup(),
			"jumpcloud_user_group_membership": resourceUserGroupMembership(),
		},
		ConfigureFunc: providerConfigure,
	}
	return p
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"api_key": "The x-api-key header used to connect to JumpCloud.",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIKey: d.Get("api_key").(string),
	}

	return config.Client()
}

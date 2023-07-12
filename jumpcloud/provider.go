package jumpcloud

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// Provider instantiates a terraform provider for Jumpcloud
// This includes all operations on all supported resources and
// global Jumpcloud parameters
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUMPCLOUD_API_KEY", nil),
				Description: descriptions["api_key"],
			},
			"org_id": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUMPCLOUD_ORG_ID", nil),
				Description: descriptions["org_id"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"jumpcloud_user":                       resourceUser(),
			"jumpcloud_user_group":                 resourceUserGroup(),
			"jumpcloud_user_group_membership":      resourceUserGroupMembership(),
			"jumpcloud_user_group_ldap_membership": resourceUserGroupLdapMembership(),
			"jumpcloud_user_group_association":     resourceUserGroupAssociation(),
			"jumpcloud_system_group":               resourceGroupsSystem(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"jumpcloud_ldap_server": dataResourceLdapServer(),
			"jumpcloud_application": dataSourceApplication(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"api_key": "The x-api-key header used to connect to JumpCloud.",
		"org_id":  "The x-org-id header used to connect to JumpCloud.",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIKey: d.Get("api_key").(string),
		OrgID:  d.Get("org_id").(string),
	}

	return config.Client()
}

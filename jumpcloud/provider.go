package jumpcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider instantiates a terraform provider for Jumpcloud
// This includes all operations on all supported resources and
// global Jumpcloud parameters
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("JUMPCLOUD_API_KEY", nil),
				Description: descriptions["api_key"],
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUMPCLOUD_CLIENT_ID", nil),
				Description: descriptions["client_id"],
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("JUMPCLOUD_CLIENT_SECRET", nil),
				Description: descriptions["client_secret"],
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
		ConfigureContextFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"api_key": "Legacy x-api-key used to connect to JumpCloud (admin-user API key). " +
			"Mutually exclusive with client_id/client_secret; OAuth takes precedence when both are set.",
		"client_id": "JumpCloud service-account OAuth client id. When set together with " +
			"client_secret, the provider mints a short-lived Bearer access token via the " +
			"client_credentials grant instead of using the legacy x-api-key.",
		"client_secret": "JumpCloud service-account OAuth client secret (paired with client_id).",
		"org_id":        "The x-org-id header used to connect to JumpCloud.",
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		APIKey:       d.Get("api_key").(string),
		ClientID:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
		OrgID:        d.Get("org_id").(string),
	}

	// client_id and client_secret are only meaningful as a pair.
	if (config.ClientID != "") != (config.ClientSecret != "") {
		return nil, diag.Errorf("client_id and client_secret must be set together")
	}

	hasOAuth := config.ClientID != "" && config.ClientSecret != ""
	if !hasOAuth && config.APIKey == "" {
		return nil, diag.Errorf("no JumpCloud credentials configured: set either " +
			"client_id + client_secret (service-account OAuth) or api_key (legacy x-api-key)")
	}

	client, err := config.Client()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return client, nil
}

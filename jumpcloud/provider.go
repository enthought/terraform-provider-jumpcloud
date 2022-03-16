package jumpcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider instantiates a terraform provider for Jumpcloud
// This includes all operations on all supported resources and
// global Jumpcloud parameters
func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
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
				"jumpcloud_user":                  resourceUser(),
				"jumpcloud_user_group":            resourceUserGroup(),
				"jumpcloud_user_group_membership": resourceUserGroupMembership(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"api_key": "The x-api-key header used to connect to JumpCloud.",
		"org_id":  "The x-org-id header used to connect to JumpCloud.",
	}
}

func new_init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

type apiClient struct {
	// Add whatever fields, client or connection info, etc. here
	// you would need to setup to communicate with the upstream
	// API.
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		// Setup a User-Agent for your API client (replace the provider name for yours):
		// userAgent := p.UserAgent("terraform-provider-scaffolding", version)
		// TODO: myClient.UserAgent = userAgent

		config := Config{
			APIKey: d.Get("api_key").(string),
			OrgID:  d.Get("org_id").(string),
		}

		return config.Client(), nil
	}
}

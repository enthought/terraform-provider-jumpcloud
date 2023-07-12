package jumpcloud

import (
	"context"
	"fmt"

	jcapiv1 "github.com/TheJumpCloud/jcapi-go/v1"
	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApplication() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApplicationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceApplicationRead(d *schema.ResourceData, m interface{}) error {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)
	applicationName, nameExists := d.GetOk("name")
	displayLabel, displayLabelExists := d.GetOk("display_label")

	if !nameExists && !displayLabelExists {
		return fmt.Errorf("either name or display_label must be provided")
	}

	applicationsResponse, _, err := client.ApplicationsApi.ApplicationsList(context.Background(), "_id, displayName, displayLabel", "", nil)

	if err != nil {
		return err
	}

	applications := applicationsResponse.Results

	for _, application := range applications {
		if (nameExists && application.DisplayName == applicationName) || (displayLabelExists && application.DisplayLabel == displayLabel) {
			d.SetId(application.Id)
			return nil
		}
	}

	return fmt.Errorf("no application found with the provided filters")
}

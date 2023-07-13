package jumpcloud

import (
	jcapiv1 "github.com/TheJumpCloud/jcapi-go/v1"
	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
)

// We receive a v2config from the TF base code but need a v1config to continue. So, we take the only
// preloaded element (the x-api-key) and populate the v1config with it.
func convertV2toV1Config(v2config *jcapiv2.Configuration) *jcapiv1.Configuration {
	configv1 := jcapiv1.NewConfiguration()
	configv1.AddDefaultHeader("x-api-key", v2config.DefaultHeader["x-api-key"])
	if v2config.DefaultHeader["x-org-id"] != "" {
		configv1.AddDefaultHeader("x-org-id", v2config.DefaultHeader["x-org-id"])
	}
	return configv1
}

package jumpcloud

import (
	jcapiv1 "github.com/TheJumpCloud/jcapi-go/v1"
	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
)

// We receive a v2config from the TF base code but need a v1config to continue. So, we carry over the
// auth header it was configured with — either the legacy x-api-key or the service-account
// `Authorization: Bearer` token — plus the optional x-org-id.
func convertV2toV1Config(v2config *jcapiv2.Configuration) *jcapiv1.Configuration {
	configv1 := jcapiv1.NewConfiguration()
	if v := v2config.DefaultHeader["Authorization"]; v != "" {
		configv1.AddDefaultHeader("Authorization", v)
	}
	if v := v2config.DefaultHeader["x-api-key"]; v != "" {
		configv1.AddDefaultHeader("x-api-key", v)
	}
	if v2config.DefaultHeader["x-org-id"] != "" {
		configv1.AddDefaultHeader("x-org-id", v2config.DefaultHeader["x-org-id"])
	}
	return configv1
}

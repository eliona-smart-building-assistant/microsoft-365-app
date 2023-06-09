/*
 * App template API
 *
 * API to access and configure the app template
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

// Configuration - Each configuration defines access to MS Graph.
type Configuration struct {

	// Internal identifier for the configured API (created automatically).
	Id *int64 `json:"id,omitempty"`

	// Client ID
	ClientId string `json:"clientId,omitempty"`

	// Client Secret
	ClientSecret string `json:"clientSecret,omitempty"`

	// Tenant ID
	TenantId string `json:"tenantId,omitempty"`

	// Username
	Username *string `json:"username,omitempty"`

	// Password
	Password *string `json:"password,omitempty"`

	// Flag to enable or disable fetching from this API
	Enable *bool `json:"enable,omitempty"`

	// Interval in seconds for collecting data from API
	RefreshInterval int32 `json:"refreshInterval,omitempty"`

	// Timeout in seconds
	RequestTimeout *int32 `json:"requestTimeout,omitempty"`

	// Array of rules combined by logical OR
	AssetFilter [][]FilterRule `json:"assetFilter,omitempty"`

	// Set to `true` by the app when running and to `false` when app is stopped
	Active *bool `json:"active,omitempty"`

	// List of Eliona project ids for which this device should collect data. For each project id all smart devices are automatically created as an asset in Eliona. The mapping to Eliona assets is stored as an asset mapping in the MS Graph app.
	ProjectIDs *[]string `json:"projectIDs,omitempty"`
}

// AssertConfigurationRequired checks if the required fields are not zero-ed
func AssertConfigurationRequired(obj Configuration) error {
	if err := AssertRecurseFilterRuleRequired(obj.AssetFilter); err != nil {
		return err
	}
	return nil
}

// AssertRecurseConfigurationRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Configuration (e.g. [][]Configuration), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseConfigurationRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aConfiguration, ok := obj.(Configuration)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertConfigurationRequired(aConfiguration)
	})
}

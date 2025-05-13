package azure

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type azureConfig struct {
	TenantID              *string  `hcl:"tenant_id"`
	SubscriptionID        *string  `hcl:"subscription_id"`
	ClientID              *string  `hcl:"client_id"`
	ClientSecret          *string  `hcl:"client_secret"`
	CertificatePath       *string  `hcl:"certificate_path"`
	CertificatePassword   *string  `hcl:"certificate_password"`
	Username              *string  `hcl:"username"`
	Password              *string  `hcl:"password"`
	Environment           *string  `hcl:"environment"`
	MaxErrorRetryAttempts *int     `hcl:"max_error_retry_attempts"`
	MinErrorRetryDelay    *int32   `hcl:"min_error_retry_delay"`
	IgnoreErrorCodes      []string `hcl:"ignore_error_codes,optional"`
	ResourceGroup         *string  `hcl:"resource_group,optional"`
	ResourceGroups        []string `hcl:"resource_groups,optional"`
}

func ConfigInstance() interface{} {
	return &azureConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) azureConfig {
	if connection == nil || connection.Config == nil {
		return azureConfig{}
	}
	config, _ := connection.Config.(azureConfig)

	// Normalize the single ResourceGroup if it's set
	if config.ResourceGroup != nil {
		// Normalize the single resource group
		normalizedRG := NormalizeResourceGroup(*config.ResourceGroup)
		config.ResourceGroup = &normalizedRG
	}

	if config.ResourceGroups != nil {
		if len(config.ResourceGroups) == 0 {
			// Empty resource_groups array means no resource group filtering
			plugin.Logger(context.Background()).Warn("connection_config", "connection_name", connection.Name, "empty_resource_groups", "resource_groups = [] means no resource group filtering will be applied")
		} else {
			// Normalize resource group names
			for i, rg := range config.ResourceGroups {
				config.ResourceGroups[i] = NormalizeResourceGroup(rg)
			}
		}
	}

	return config
}

func NormalizeResourceGroup(resourceGroup string) string {
	// ensure resource groups are lower case, to work consistently in matching
	// and comparisons
	return strings.ToLower(resourceGroup)
}

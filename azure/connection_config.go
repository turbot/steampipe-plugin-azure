package azure

import (
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
	MaxErrorRetryAttempts *int64   `hcl:"max_error_retry_attempts"`
	MinErrorRetryDelay    *int64   `hcl:"min_error_retry_delay"`
	IgnoreErrorCodes      []string `hcl:"ignore_error_codes,optional"`
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
	return config
}

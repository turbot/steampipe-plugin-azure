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
	MaxErrorRetryAttempts *int     `hcl:"max_error_retry_attempts"`
	MinErrorRetryDelay    *int32   `hcl:"min_error_retry_delay"`
	IgnoreErrorCodes      []string `hcl:"ignore_error_codes,optional"`

	// Storage data-plane auth configuration
	AuthMode               *string `hcl:"auth_mode"`                 // aad (default) | shared_key | sas | auto (deprecated)
	StorageAccountKey      *string `hcl:"storage_account_key"`       // explicit shared key
	StorageSASToken        *string `hcl:"storage_sas_token"`         // SAS token (with or without leading '?')
	AllowStorageKeyListing *bool   `hcl:"allow_storage_key_listing"` // permit ListKeys call when using shared_key without explicit key
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

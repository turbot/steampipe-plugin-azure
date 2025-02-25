package azure

import (
	"context"
	"reflect"
	"time"

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

type RetryRule struct {
	MaxErrorRetryAttempts *int
	MinErrorRetryDelay    *time.Duration
}

// Customize the RetryRules to implement exponential backoff retry
func getRetryRules(connection *plugin.Connection) *RetryRule {
	connectionConfig := GetConfig(connection)

	if connectionConfig.MaxErrorRetryAttempts != nil && *connectionConfig.MaxErrorRetryAttempts < 1 {
		panic("connection config has invalid value for \"max_error_retry_attempts\", it must be greater than or equal to 1")
	}

	if connectionConfig.MinErrorRetryDelay != nil && *connectionConfig.MinErrorRetryDelay < 1 {
		panic("connection config has invalid value for \"min_error_retry_delay\", it must be greater than or equal to 1")
	}

	// Fallback to SDK default value
	// https://github.com/Azure/go-autorest/blob/main/autorest/client.go#L39
	maxRetries := 3
	minDelay := 30 * time.Second

	if connectionConfig.MaxErrorRetryAttempts != nil {
		maxRetries = int(*connectionConfig.MaxErrorRetryAttempts)
	}

	if connectionConfig.MinErrorRetryDelay != nil {
		minDelay = time.Duration(*connectionConfig.MinErrorRetryDelay) * time.Millisecond
	}

	return &RetryRule{
		MaxErrorRetryAttempts: &maxRetries,
		MinErrorRetryDelay:    &minDelay,
	}
}

// ApplyRetryRules applies retry settings to any Azure SDK client
func ApplyRetryRules(ctx context.Context, client interface{}, connection *plugin.Connection) {
	v := reflect.ValueOf(client).Elem()

	retryRules := getRetryRules(connection)

	// Set RetryAttempts if the field exists
	if field := v.FieldByName("RetryAttempts"); field.IsValid() && field.CanSet() {
		retryAttempts := int64(*retryRules.MaxErrorRetryAttempts)
		field.SetInt(retryAttempts)
	} else if field := v.FieldByName("RetryAttempts"); !field.IsValid() || !field.CanSet() {
		plugin.Logger(ctx).Warn("'RetryAttempts' could not be set")
	}

	// Set RetryDuration if the field exists
	if field := v.FieldByName("RetryDuration"); field.IsValid() && field.CanSet() {
		field.SetInt(int64(retryRules.MinErrorRetryDelay.Milliseconds()))
	} else if field := v.FieldByName("RetryDuration"); !field.IsValid() || !field.CanSet() {
		plugin.Logger(ctx).Warn("'RetryDuration' could not be set")
	}
}

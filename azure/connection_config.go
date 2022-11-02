package azure

import (
	"math"
	"math/rand"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/schema"
)

type azureConfig struct {
	TenantID              *string `cty:"tenant_id"`
	SubscriptionID        *string `cty:"subscription_id"`
	ClientID              *string `cty:"client_id"`
	ClientSecret          *string `cty:"client_secret"`
	CertificatePath       *string `cty:"certificate_path"`
	CertificatePassword   *string `cty:"certificate_password"`
	Username              *string `cty:"username"`
	Password              *string `cty:"password"`
	Environment           *string `cty:"environment"`
	MaxErrorRetryAttempts *int    `cty:"max_error_retry_attempts"`
	MinErrorRetryDelay    *int    `cty:"min_error_retry_delay"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"tenant_id": {
		Type: schema.TypeString,
	},
	"subscription_id": {
		Type: schema.TypeString,
	},
	"client_id": {
		Type: schema.TypeString,
	},
	"client_secret": {
		Type: schema.TypeString,
	},
	"certificate_path": {
		Type: schema.TypeString,
	},
	"certificate_password": {
		Type: schema.TypeString,
	},
	"username": {
		Type: schema.TypeString,
	},
	"password": {
		Type: schema.TypeString,
	},
	"environment": {
		Type: schema.TypeString,
	},
	"max_error_retry_attempts": {
		Type: schema.TypeInt,
	},
	"min_error_retry_delay": {
		Type: schema.TypeInt,
	},
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
	MinErrorRetryDelay    time.Duration
}

// Customize the RetryRules to implement exponential backoff retry
func getRetryRules(connection *plugin.Connection, retryCount int) *RetryRule {
	connectionConfig := GetConfig(connection)

	maxRetries := 9
	var minDelay time.Duration = 25 * time.Millisecond

	if connectionConfig.MaxErrorRetryAttempts != nil && *connectionConfig.MaxErrorRetryAttempts > 9 {
		maxRetries = *connectionConfig.MaxErrorRetryAttempts
	}

	if connectionConfig.MinErrorRetryDelay != nil && *connectionConfig.MinErrorRetryDelay > 25 {
		minDelay = time.Duration(*connectionConfig.MinErrorRetryDelay) * time.Millisecond
	}

	// If errors are caused by load, retries can be ineffective if all API request retry at the same time.
	// To avoid this problem added a jitter of "+/-20%" with delay time.
	// For example, if the delay is 25ms, the final delay could be between 20 and 30ms.
	var jitter = float64(rand.Intn(120-80)+80) / 100

	// Creates a new exponential backoff using the starting value of
	// minDelay and (minDelay * 3^retrycount) * jitter on each failure
	// For example, with a min delay time of 25ms: 23.25ms, 63ms, 238.5ms, 607.4ms, 2s, 5.22s, 20.31s..., up to max.
	retryTime := time.Duration(int(float64(int(minDelay.Nanoseconds())*int(math.Pow(3, float64(retryCount)))) * jitter))

	// Cap retry time at 5 minuets to avoid too long a wait
	if retryTime > time.Duration(5*time.Minute) {
		retryTime = time.Duration(5 * time.Minute)
	}

	return &RetryRule{
		MaxErrorRetryAttempts: &maxRetries,
		MinErrorRetryDelay:    retryTime,
	}
}

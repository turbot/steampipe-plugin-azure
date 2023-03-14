package azure

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type azureConfig struct {
	TenantID            *string  `cty:"tenant_id"`
	SubscriptionID      *string  `cty:"subscription_id"`
	ClientID            *string  `cty:"client_id"`
	ClientSecret        *string  `cty:"client_secret"`
	CertificatePath     *string  `cty:"certificate_path"`
	CertificatePassword *string  `cty:"certificate_password"`
	Username            *string  `cty:"username"`
	Password            *string  `cty:"password"`
	Environment         *string  `cty:"environment"`
	IgnoreErrorCodes    []string `cty:"ignore_error_codes"`
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
	"ignore_error_codes": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
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

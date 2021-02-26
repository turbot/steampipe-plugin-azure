package azure

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type azureConfig struct {
	TenantID            *string `cty:"tenant_id"`
	SubscriptionID      *string `cty:"subscription_id"`
	ClientID            *string `cty:"client_id"`
	ClientSecret        *string `cty:"client_secret"`
	CertificatePath     *string `cty:"client_certificate_path"`
	CertificatePassword *string `cty:"client_certificate_password"`
	Username            *string `cty:"username"`
	Password            *string `cty:"password"`
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
	"client_certificate_path": {
		Type: schema.TypeString,
	},
	"client_certificate_password": {
		Type: schema.TypeString,
	},
	"username": {
		Type: schema.TypeString,
	},
	"password": {
		Type: schema.TypeString,
	},
	// "use_msi": {
	// 	Type: schema.TypeString,
	// },
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

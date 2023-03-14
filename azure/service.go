package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/azure/cli"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Session info
type Session struct {
	Authorizer              autorest.Authorizer
	CloudEnvironment        string
	Expires                 *time.Time
	GraphEndpoint           string
	ResourceManagerEndpoint string
	StorageEndpointSuffix   string
	SubscriptionID          string
	TenantID                string
}

/*
	GetNewSession creates an session configured from (~/.steampipe/config, environment variables and CLI) in the order:

1. Client secret
2. Client certificate
3. Username and password
4. MSI
5. CLI
*/
func GetNewSession(ctx context.Context, d *plugin.QueryData, tokenAudience string) (session *Session, err error) {
	logger := plugin.Logger(ctx)

	cacheKey := "GetNewSession" + tokenAudience
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		session = cachedData.(*Session)
		if session.Expires != nil && WillExpireIn(*session.Expires, 0) {
			logger.Trace("GetNewSession", "cache expired", "delete cache and obtain new session token")
			d.ConnectionManager.Cache.Delete(cacheKey)
		} else {
			return cachedData.(*Session), nil
		}
	}

	logger.Debug("Auth session not found in cache, creating new session")

	var subscriptionID, tenantID string
	settings := auth.EnvironmentSettings{
		Values:      map[string]string{},
		Environment: azure.PublicCloud, // Set public cloud as default
	}

	azureConfig := GetConfig(d.Connection)

	if azureConfig.TenantID != nil {
		tenantID = *azureConfig.TenantID
		settings.Values[auth.TenantID] = *azureConfig.TenantID
	} else {
		tenantID = os.Getenv(auth.TenantID)
		settings.Values[auth.TenantID] = os.Getenv(auth.TenantID)
	}

	if azureConfig.SubscriptionID != nil {
		subscriptionID = *azureConfig.SubscriptionID
		settings.Values[auth.SubscriptionID] = *azureConfig.SubscriptionID
	} else {
		subscriptionID = os.Getenv(auth.SubscriptionID)
		settings.Values[auth.SubscriptionID] = os.Getenv(auth.SubscriptionID)
	}

	if azureConfig.ClientID != nil {
		settings.Values[auth.ClientID] = *azureConfig.ClientID
	} else {
		settings.Values[auth.ClientID] = os.Getenv(auth.ClientID)
	}

	if azureConfig.ClientSecret != nil {
		settings.Values[auth.ClientSecret] = *azureConfig.ClientSecret
	} else {
		settings.Values[auth.ClientSecret] = os.Getenv(auth.ClientSecret)
	}

	if azureConfig.CertificatePath != nil {
		settings.Values[auth.CertificatePath] = *azureConfig.CertificatePath
	} else {
		settings.Values[auth.CertificatePath] = os.Getenv(auth.CertificatePath)
	}

	if azureConfig.CertificatePassword != nil {
		settings.Values[auth.CertificatePassword] = *azureConfig.CertificatePassword
	} else {
		settings.Values[auth.CertificatePassword] = os.Getenv(auth.CertificatePassword)
	}

	if azureConfig.Username != nil {
		settings.Values[auth.Username] = *azureConfig.Username
	} else {
		settings.Values[auth.Username] = os.Getenv(auth.Username)
	}

	if azureConfig.Password != nil {
		settings.Values[auth.Password] = *azureConfig.Password
	} else {
		settings.Values[auth.Password] = os.Getenv(auth.Password)
	}

	if azureConfig.Environment != nil {
		env, err := azure.EnvironmentFromName(*azureConfig.Environment)
		if err != nil {
			logger.Error("GetNewSession", "Error getting environment from name with config environment", err)
			return nil, err
		}
		settings.Environment = env
		settings.Values[auth.EnvironmentName] = *azureConfig.Environment
	} else {
		env := azure.PublicCloud
		envName, ok := os.LookupEnv(auth.EnvironmentName)
		if ok {
			env, err = azure.EnvironmentFromName(envName)
			if err != nil {
				logger.Error("GetNewSession", "Error getting environment from name with no config environment", err)
				return nil, err
			}
			settings.Values[auth.EnvironmentName] = envName
		}
		settings.Environment = env
	}

	authMethod, resource, err := getApplicableAuthorizationDetails(ctx, settings, tokenAudience)
	if err != nil {
		logger.Error("GetNewSession", "getApplicableAuthorizationDetails error", err)
		return nil, err
	}
	settings.Values[auth.Resource] = resource

	var authorizer autorest.Authorizer
	var expiresOn *time.Time

	// so if it was not in cache - create session
	switch authMethod {
	case "Environment":
		logger.Trace("Creating new session authorizer from environment")
		authorizer, err = settings.GetAuthorizer()
		if err != nil {
			logger.Error("GetNewSession", "NewAuthorizerFromEnvironmentWithResource error", err)
			return nil, err
		}

	// Get the subscription ID and tenant ID for "GRAPH" token audience
	case "CLI":
		logger.Trace("Getting token for authorizer from Azure CLI")
		token, err := cli.GetTokenFromCLI(resource)
		if err != nil {
			logger.Error("GetNewSession", "get_token_from_cli_error", err)
			return nil, err
		}

		adalToken, err := token.ToADALToken()
		expiresOn = types.Time(adalToken.Expires())
		logger.Trace("GetNewSession", "Getting token for authorizer from Azure CLI, expiresOn", expiresOn.Local())

		if err != nil {
			logger.Error("GetNewSession", "Get token from Azure CLI error", err)
			// Check if the password was changed and the session token is stored in the system, or if the CLI is outdated
			if strings.Contains(err.Error(), "invalid_grant") {
				return nil, fmt.Errorf("ValidationError: The credential data used by the CLI has expired because you might have changed or reset the password. Please clear your browser's cookies and run 'az login'.")
			}
			return nil, err
		}
		authorizer = autorest.NewBearerAuthorizer(&adalToken)
	default:
		return nil, fmt.Errorf("invalid Azure authentication method: %s", authMethod)
	}

	// Get the subscription ID and tenant ID from CLI if not set in connection
	// config or environment variables
	if authMethod == "CLI" && (settings.Values[auth.SubscriptionID] == "" || settings.Values[auth.TenantID] == "") {
		logger.Trace("Getting subscription ID and/or tenant ID from from Azure CLI")
		subscription, err := getSubscriptionFromCLI(resource)
		if err != nil {
			logger.Error("GetNewSession", "getSubscriptionFromCLI error", err)
			return nil, err
		}
		tenantID = subscription.TenantID

		// Subscription ID set in config file or environment variable takes
		// precedence over the subscription ID set in the CLI
		if subscriptionID == "" {
			subscriptionID = subscription.SubscriptionID
			logger.Trace("Setting subscription ID from Azure CLI", "subscription_id", subscriptionID)
		}
	}

	sess := &Session{
		Authorizer:              authorizer,
		CloudEnvironment:        settings.Environment.Name,
		Expires:                 expiresOn,
		GraphEndpoint:           settings.Environment.GraphEndpoint,
		ResourceManagerEndpoint: settings.Environment.ResourceManagerEndpoint,
		StorageEndpointSuffix:   settings.Environment.StorageEndpointSuffix,
		SubscriptionID:          subscriptionID,
		TenantID:                tenantID,
	}

	var expireMins time.Duration
	if expiresOn != nil {
		expireMins = time.Until(*sess.Expires)
	} else {
		// Cache for 55 minutes to avoid expiry issue
		expireMins = time.Minute * 55
	}

	logger.Debug("Session saved in cache", "expiration_time", expireMins)
	d.ConnectionManager.Cache.SetWithTTL(cacheKey, sess, expireMins)

	return sess, err
}

func getApplicableAuthorizationDetails(ctx context.Context, settings auth.EnvironmentSettings, tokenAudience string) (authMethod string, resource string, err error) {
	logger := plugin.Logger(ctx)
	subscriptionID := settings.Values[auth.SubscriptionID]
	tenantID := settings.Values[auth.TenantID]
	clientID := settings.Values[auth.ClientID]
	// Azure environment name
	environmentName := settings.Values[auth.EnvironmentName]

	// CLI is the default authentication method
	authMethod = "CLI"
	if subscriptionID == "" || (subscriptionID == "" && tenantID == "") {
		authMethod = "CLI"
	} else if subscriptionID != "" && tenantID != "" && clientID != "" {
		// Works for client secret credentials, client certificate credentials, resource owner password, and managed identities
		authMethod = "Environment"
	}

	logger.Debug("getApplicableAuthorizationDetails", "auth_method", authMethod)

	var environment azure.Environment
	// Get the environment endpoint to be used for authorization
	if environmentName == "" {
		settings.Environment = azure.PublicCloud
	} else {
		environment, err = azure.EnvironmentFromName(environmentName)
		if err != nil {
			logger.Error("getApplicableAuthorizationDetails", "get_environment_name_error", err)
			return
		}
		settings.Environment = environment
	}

	logger.Debug("getApplicableAuthorizationDetails", "token_audience", tokenAudience)

	switch tokenAudience {
	case "GRAPH":
		resource = settings.Environment.GraphEndpoint
	case "VAULT":
		resource = strings.TrimSuffix(settings.Environment.KeyVaultEndpoint, "/")
	case "MANAGEMENT":
		resource = settings.Environment.ResourceManagerEndpoint
	default:
		resource = settings.Environment.ResourceManagerEndpoint
	}

	logger.Debug("getApplicableAuthorizationDetails", "resource", resource)

	return
}

type subscription struct {
	SubscriptionID string `json:"subscriptionID,omitempty"`
	TenantID       string `json:"tenantID,omitempty"`
}

// https://github.com/Azure/go-autorest/blob/3fb5326fea196cd5af02cf105ca246a0fba59021/autorest/azure/cli/token.go#L126
// NewAuthorizerFromCLIWithResource creates an Authorizer configured from Azure CLI 2.0 for local development scenarios.
func getSubscriptionFromCLI(resource string) (*subscription, error) {
	// This is the path that a developer can set to tell this class what the install path for Azure CLI is.
	const azureCLIPath = "AzureCLIPath"

	// The default install paths are used to find Azure CLI. This is for security, so that any path in the calling program's Path environment is not used to execute Azure CLI.
	azureCLIDefaultPathWindows := fmt.Sprintf("%s\\Microsoft SDKs\\Azure\\CLI2\\wbin; %s\\Microsoft SDKs\\Azure\\CLI2\\wbin", os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramFiles"))

	// Default path for non-Windows.
	const azureCLIDefaultPath = "/bin:/sbin:/usr/bin:/usr/local/bin"

	// Validate resource, since it gets sent as a command line argument to Azure CLI
	const invalidResourceErrorTemplate = "Resource %s is not in expected format. Only alphanumeric characters, [dot], [colon], [hyphen], and [forward slash] are allowed."
	match, err := regexp.MatchString("^[0-9a-zA-Z-.:/]+$", resource)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, fmt.Errorf(invalidResourceErrorTemplate, resource)
	}

	// Execute Azure CLI to get token
	var cliCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cliCmd = exec.Command(fmt.Sprintf("%s\\system32\\cmd.exe", os.Getenv("windir")))
		cliCmd.Env = os.Environ()
		cliCmd.Env = append(cliCmd.Env, fmt.Sprintf("PATH=%s;%s", os.Getenv(azureCLIPath), azureCLIDefaultPathWindows))
		cliCmd.Args = append(cliCmd.Args, "/c", "az")
	} else {
		cliCmd = exec.Command("az")
		cliCmd.Env = os.Environ()
		cliCmd.Env = append(cliCmd.Env, fmt.Sprintf("PATH=%s:%s", os.Getenv(azureCLIPath), azureCLIDefaultPath))
	}
	cliCmd.Args = append(cliCmd.Args, "account", "get-access-token", "-o", "json", "--resource", resource)

	var stderr bytes.Buffer
	cliCmd.Stderr = &stderr

	output, err := cliCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Invoking Azure CLI failed with the following error: %v", err)
	}

	var tokenResponse map[string]interface{}
	err = json.Unmarshal(output, &tokenResponse)
	if err != nil {
		return nil, err
	}

	return &subscription{
		SubscriptionID: tokenResponse["subscription"].(string),
		TenantID:       tokenResponse["tenant"].(string),
	}, nil
}

// WillExpireIn returns true if the Token will expire after the passed time.Duration interval
// from now, false otherwise.
func WillExpireIn(t time.Time, d time.Duration) bool {
	return !t.After(time.Now().Add(d))
}

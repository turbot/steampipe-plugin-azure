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
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// Session info
type Session struct {
	SubscriptionID          string
	TenantID                string
	Authorizer              autorest.Authorizer
	Expires                 *time.Time
	ResourceManagerEndpoint string
	StorageEndpointSuffix   string
}

/* GetNewSession creates an session configured from (~/.steampipe/config, environment variables and CLI) in the order:
1. Client secret
2. Client certificate
3. Username and password
4. MSI
5. CLI
*/
func GetNewSession(ctx context.Context, d *plugin.QueryData, tokenAudience string) (session *Session, err error) {
	logger := plugin.Logger(ctx)
	cacheKey := "GetNewSession"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		session = cachedData.(*Session)
		if session.Expires != nil && WillExpireIn(*session.Expires, 0) {
			d.ConnectionManager.Cache.Delete("GetNewSession")
		} else {
			return cachedData.(*Session), nil
		}
	}

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
			logger.Debug("GetNewSession_", "Error getting environment from name with config environment", err)
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
				logger.Debug("GetNewSession_", "Error getting environment from name with no config environment", err)
				return nil, err
			}
			settings.Values[auth.EnvironmentName] = envName
		}
		settings.Environment = env
	}

	authMethod, resource, err := getApplicableAuthorizationDetails(ctx, settings, tokenAudience)
	if err != nil {
		logger.Debug("GetNewSession__", "getApplicableAuthorizationDetails error", err)
		return nil, err
	}
	settings.Values[auth.Resource] = resource

	var authorizer autorest.Authorizer
	var expiresOn time.Time

	// so if it was not in cache - create session
	switch authMethod {
	case "Environment":
		authorizer, err = settings.GetAuthorizer()
		if err != nil {
			logger.Debug("GetNewSession__", "NewAuthorizerFromEnvironmentWithResource error", err)
			return nil, err
		}

	// In this case need get the details of SUBSCRIPTION_ID and TENANT_ID if
	// tokenAudience is GRAPH
	case "CLI":
		authorizer, err = auth.NewAuthorizerFromCLIWithResource(resource)
		if err != nil {
			logger.Debug("GetNewSession__", "NewAuthorizerFromCLIWithResource error", err)

			// Check if the password was changed and the session token is stored in
			// the system, or if the CLI is outdated
			if strings.Contains(err.Error(), "invalid_grant") {
				return nil, fmt.Errorf("ValidationError: The credential data used by the CLI has expired because you might have changed or reset the password. Please clear your browser's cookies and run 'az login'.")
			}
			return nil, err
		}
	default:
		token, err := cli.GetTokenFromCLI(resource)
		if err != nil {
			return nil, err
		}

		adalToken, err := token.ToADALToken()
		expiresOn = adalToken.Expires()

		if err != nil {
			logger.Debug("GetNewSession__", "NewAuthorizerFromCLIWithResource error", err)

			if strings.Contains(err.Error(), "invalid_grant") {
				return nil, fmt.Errorf("ValidationError: The credential data used by the CLI has expired because you might have changed or reset the password. Please clear your browser's cookies and run 'az login'.")
			}
			return nil, err
		}
		authorizer = autorest.NewBearerAuthorizer(&adalToken)
	}

	if authMethod == "CLI" {
		subscription, err := getSubscriptionFromCLI(resource)
		if err != nil {
			logger.Debug("GetNewSession__", "getSubscriptionFromCLI error", err)
			return nil, err
		}
		tenantID = subscription.TenantID

		// If "AZURE_SUBSCRIPTION_ID" is set then it will take precedence over the subscription set in the CLI
		if subscriptionID == "" {
			subscriptionID = subscription.SubscriptionID
		}
	}

	sess := &Session{
		SubscriptionID:          subscriptionID,
		Authorizer:              authorizer,
		TenantID:                tenantID,
		Expires:                 &expiresOn,
		ResourceManagerEndpoint: settings.Environment.ResourceManagerEndpoint,
		StorageEndpointSuffix:   settings.Environment.StorageEndpointSuffix,
	}

	if sess.Expires != nil {
		d.ConnectionManager.Cache.SetWithTTL(cacheKey, sess, time.Until(*sess.Expires))
	} else {
		// Cache for 55 minutes to avoid expiry issue
		d.ConnectionManager.Cache.SetWithTTL(cacheKey, sess, time.Minute*55)
	}

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

	logger.Trace("getApplicableAuthorizationDetails_", "Auth Method: ", authMethod)

	var environment azure.Environment
	// Get the environment endpoint to be used for authorization
	if environmentName == "" {
		settings.Environment = azure.PublicCloud
	} else {
		environment, err = azure.EnvironmentFromName(environmentName)
		if err != nil {
			logger.Error("Unable to get azure environment", "ERROR", err)
			return
		}
		settings.Environment = environment
	}

	logger.Trace("getApplicableAuthorizationDetails_", "tokenAudience: ", tokenAudience)

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

	logger.Trace("getApplicableAuthorizationDetails_", "resource: ", resource)

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

package azure

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	armstorage "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	azblob "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	azqueue "github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// buildBlobServiceClient creates an azblob Client honoring modes: auto|aad|shared_key with reusable fallback logic.
func buildBlobServiceClient(ctx context.Context, d *plugin.QueryData, tokenCred azcore.TokenCredential, authMode, accountName, resourceGroup, subscriptionID, storageEndpointSuffix string, allowShared bool) (*azblob.Client, error) {
	endpoint := fmt.Sprintf("https://%s.blob.%s", accountName, storageEndpointSuffix)
	aadClient, err := azblob.NewClient(endpoint, tokenCred, nil)
	if err != nil {
		return nil, err
	}
	mode := normalizeDataPlaneMode(authMode)
	// Shared key builder closure
	buildShared := func() (*azblob.Client, error) {
		return buildBlobServiceClientSharedKey(ctx, tokenCred, accountName, resourceGroup, subscriptionID, storageEndpointSuffix)
	}
	probe := func() error {
		one := int32(1)
		pager := aadClient.NewListContainersPager(&azblob.ListContainersOptions{MaxResults: &one})
		if pager.More() {
			_, perr := pager.NextPage(ctx)
			return perr
		}
		return nil
	}
	return attemptAADWithFallback(ctx, mode, allowShared, "azure_storage_blob", buildShared, probe, aadClient)
}

// buildBlobServiceClientSharedKey returns a shared key client by listing account keys (used for explicit shared_key or auto fallback)
func buildBlobServiceClientSharedKey(ctx context.Context, tokenCred azcore.TokenCredential, accountName, resourceGroup, subscriptionID, storageEndpointSuffix string) (*azblob.Client, error) {
	endpoint := fmt.Sprintf("https://%s.blob.%s", accountName, storageEndpointSuffix)
	key, err := getFirstStorageAccountKey(ctx, tokenCred, subscriptionID, resourceGroup, accountName)
	if err != nil {
		return nil, err
	}
	cred, err := azblob.NewSharedKeyCredential(accountName, key)
	if err != nil {
		return nil, err
	}
	return azblob.NewClientWithSharedKeyCredential(endpoint, cred, nil)
}

// translateBlobError provides friendly error messages for common storage authorization cases.
func translateBlobError(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if strings.Contains(msg, "AuthorizationPermissionMismatch") || strings.Contains(strings.ToLower(msg), "authorizationfailure") {
		return errors.New("authorization failed (possible missing Storage Blob Data Reader role)")
	}
	return err
}

// buildQueueServiceClient creates a track2 azqueue Client for queue service (account-level) per auth_mode.
func buildQueueServiceClient(ctx context.Context, d *plugin.QueryData, tokenCred azcore.TokenCredential, authMode, accountName, resourceGroup, subscriptionID, storageEndpointSuffix string, allowShared bool) (*azqueue.ServiceClient, error) {
	endpoint := fmt.Sprintf("https://%s.queue.%s", accountName, storageEndpointSuffix)
	aadClient, err := azqueue.NewServiceClient(endpoint, tokenCred, nil)
	if err != nil {
		return nil, err
	}
	mode := normalizeDataPlaneMode(authMode)
	buildShared := func() (*azqueue.ServiceClient, error) {
		return buildQueueServiceClientSharedKey(ctx, tokenCred, accountName, resourceGroup, subscriptionID, storageEndpointSuffix)
	}
	probe := func() error {
		one := int32(1)
		pager := aadClient.NewListQueuesPager(&azqueue.ListQueuesOptions{MaxResults: &one})
		if pager.More() {
			_, perr := pager.NextPage(ctx)
			return perr
		}
		return nil
	}
	return attemptAADWithFallback(ctx, mode, allowShared, "azure_storage_queue", buildShared, probe, aadClient)
}

func buildQueueServiceClientSharedKey(ctx context.Context, tokenCred azcore.TokenCredential, accountName, resourceGroup, subscriptionID, storageEndpointSuffix string) (*azqueue.ServiceClient, error) {
	endpoint := fmt.Sprintf("https://%s.queue.%s", accountName, storageEndpointSuffix)
	key, err := getFirstStorageAccountKey(ctx, tokenCred, subscriptionID, resourceGroup, accountName)
	if err != nil {
		return nil, err
	}
	cred, err := azqueue.NewSharedKeyCredential(accountName, key)
	if err != nil {
		return nil, err
	}
	return azqueue.NewServiceClientWithSharedKeyCredential(endpoint, cred, nil)
}

// helper to detect authorization failure eligible for auto fallback
func isAuthFailure(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "authorization") || strings.Contains(msg, "permission") || strings.Contains(msg, "auth failed") || strings.Contains(msg, "403")
}

// normalizeDataPlaneMode maps empty string to auto and lowercases value.
func normalizeDataPlaneMode(m string) string {
	if m == "" {
		return "auto"
	}
	return strings.ToLower(m)
}

// getFirstStorageAccountKey returns the first key value for the given storage account.
func getFirstStorageAccountKey(ctx context.Context, tokenCred azcore.TokenCredential, subscriptionID, resourceGroup, accountName string) (string, error) {
	acctClient, err := armstorage.NewAccountsClient(subscriptionID, tokenCred, nil)
	if err != nil {
		return "", err
	}
	keys, err := acctClient.ListKeys(ctx, resourceGroup, accountName, nil)
	if err != nil {
		return "", err
	}
	if len(keys.Keys) == 0 || keys.Keys[0].Value == nil {
		return "", fmt.Errorf("no storage account keys returned for '%s'", accountName)
	}
	return *keys.Keys[0].Value, nil
}

// attemptAADWithFallback encapsulates AAD->shared_key selection for auto mode.
func attemptAADWithFallback[T any](ctx context.Context, mode string, allowShared bool, logComponent string, buildShared func() (T, error), probeAAD func() error, aadClient T) (T, error) {
	switch mode {
	case "aad":
		return aadClient, nil
	case "shared_key":
		if !allowShared {
			var zero T
			return zero, fmt.Errorf("shared key access disabled on storage account")
		}
		return buildShared()
	case "auto":
		if err := probeAAD(); err != nil {
			if isAuthFailure(err) && allowShared {
				sharedClient, skErr := buildShared()
				if skErr == nil {
					plugin.Logger(ctx).Debug(logComponent, "auth_fallback", "using shared key after AAD denial")
					return sharedClient, nil
				}
				plugin.Logger(ctx).Warn(logComponent, "shared_key_fallback_failed", skErr.Error())
			}
		}
		return aadClient, nil
	default:
		plugin.Logger(ctx).Warn(logComponent, "unsupported_auth_mode", mode)
		return aadClient, nil
	}
}

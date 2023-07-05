package azure

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// isNotFoundError:: function which returns an ErrorPredicate for Azure API calls
func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		azureConfig := GetConfig(d.Connection)

		// If the get or list hydrate functions have an overriding IgnoreConfig
		// defined using the isNotFoundError function, then it should
		// also check for errors in the "ignore_error_codes" config argument
		allErrors := append(notFoundErrors, azureConfig.IgnoreErrorCodes...)
		// Added to support regex in not found errors
		for _, pattern := range allErrors {
			if strings.Contains(err.Error(), pattern) {
				return true
			}
		}
		return false
	}
}

// shouldIgnoreErrorPluginDefault:: Plugin level default function to ignore a set errors for hydrate functions based on "ignore_error_codes" config argument
func shouldIgnoreErrorPluginDefault() plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		if !hasIgnoredErrorCodes(d.Connection) {
			return false
		}

		azureConfig := GetConfig(d.Connection)
		// Added to support regex in ignoring errors
		for _, pattern := range azureConfig.IgnoreErrorCodes {
			if strings.Contains(err.Error(), pattern) {
				return true
			}
		}
		return false
	}
}

func hasIgnoredErrorCodes(connection *plugin.Connection) bool {
	azureConfig := GetConfig(connection)
	return len(azureConfig.IgnoreErrorCodes) > 0
}

func shouldRetryError(retryErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		plugin.Logger(ctx).Error("err***", err.Error())
		for _, pattern := range retryErrors {
			// handle retry error
			if strings.Contains(err.Error(), pattern) {
				return true
			}
		}
		return false
	}
}

func getDynamicRetryConfig(ctx context.Context, retryErrors []string) func(ctx context.Context, d *plugin.QueryData) *plugin.RetryConfig {
	plugin.Logger(ctx).Error("getDynamicRetryConfig")
	return func(ctx context.Context, d *plugin.QueryData) *plugin.RetryConfig {
		azureConfig := GetConfig(d.Connection)

		// set default retry limits
		retryAttempts := int64(2)
		retryDuration := int64(5)

		if azureConfig.MaxErrorRetryAttempts != nil {
			retryAttempts = *azureConfig.MaxErrorRetryAttempts
		}
		if azureConfig.MinErrorRetryDelay != nil {
			retryDuration = *azureConfig.MinErrorRetryDelay
		}
		plugin.Logger(ctx).Error("retryAttempts", retryAttempts, retryDuration)
		retryConfig := &plugin.RetryConfig{
			ShouldRetryErrorFunc: shouldRetryError(retryErrors),
			MaxAttempts:          retryAttempts,
			RetryInterval:        retryDuration,
			BackoffAlgorithm:     "Exponential",
		}

		return retryConfig
	}
}

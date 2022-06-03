package azure

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

// isNotFoundError:: function which returns an ErrorPredicate for AWS API calls
func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		azureConfig := GetConfig(d.Connection)

		// If the get or list hydrate functions have an overriding IgnoreConfig
		// defined using the isNotFoundError function, then it should
		// also check for errors in the "ignore_error_codes" config argument
		allErrors := append(notFoundErrors, azureConfig.IgnoreErrorCodes...)
		// Added to support regex in not found errors
		plugin.Logger(ctx).Error("Error Code =====>>", err.Error())
		plugin.Logger(ctx).Error("All errors =====>>", strings.Join(allErrors, ","))
		for _, pattern := range allErrors {
			if strings.Contains(err.Error(), pattern) {
				return true
			}
			// if ok, _ := path.Match(pattern, err.Error()); ok {
			// 	plugin.Logger(ctx).Error("Path Match =====>>", ok)
			// 	return true
			// }
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
			// if ok, _ := path.Match(pattern, err.Error()); ok {
			// 	return true
			// }
		}
		return false
	}
}

func hasIgnoredErrorCodes(connection *plugin.Connection) bool {
	azureConfig := GetConfig(connection)
	return len(azureConfig.IgnoreErrorCodes) > 0
}

package azure

import (
	"context"
	"path"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const matrixKeyResourceGroup = "resource_group"

// Return a matrix of resource groups for tables that support resource group filtering
// This will filter based on the resource_groups or resource_group configuration
func ResourceGroupMatrixFilter(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {
	// Get the resource groups defined in the connection configuration
	queryResourceGroups, err := listQueryResourceGroupsForConnection(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ResourceGroupMatrixFilter", "connection_name", d.Connection.Name, "query_resource_groups_error", err)
		return []map[string]interface{}{}
	}

	// Create a matrix from the resource groups
	matrix := []map[string]interface{}{}
	for _, resourceGroup := range queryResourceGroups {
		obj := map[string]interface{}{matrixKeyResourceGroup: resourceGroup}
		matrix = append(matrix, obj)
	}

	plugin.Logger(ctx).Debug("ResourceGroupMatrixFilter", "connection_name", d.Connection.Name, "matrix", matrix)
	return matrix
}

// Calculate the resource groups that the user has requested to query for this connection.
// This function supports wildcards "*" and "?" in the connection config for resource_groups.
//
// Scenarios:
// 1. When no resource_groups mentioned in connection config or resource_groups is empty:
//   - Return all the resource groups
//
// 2. When resource_groups has specific values:
//   - Return the list as-is but after verifying that the resource groups specified by the user exist
//
// 3. When resource_groups has wildcards:
//   - resource_groups = ["*"] means all resource groups
//   - resource_groups = ["prod-*"] means all resource groups starting with "prod-"
func listQueryResourceGroupsForConnection(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	// Retrieve resource groups list from the plugin connection config
	azureConfig := GetConfig(d.Connection)

	// Get all resource groups
	resourceGroups, err := getAllResourceGroups(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listQueryResourceGroupsForConnection", "connection_name", d.Connection.Name, "error", err)
		return nil, err
	}

	// If there are no resource groups defined in the config or the array is empty, return all the resource groups
	if azureConfig.ResourceGroups == nil || len(azureConfig.ResourceGroups) == 0 {
		plugin.Logger(ctx).Debug("listQueryResourceGroupsForConnection", "connection_name", d.Connection.Name, "no resource group filters configured, returning all resource groups")
		return resourceGroups, nil
	}

	// Filter to resource groups that match the patterns in the config
	var targetResourceGroups []string
	for _, pattern := range azureConfig.ResourceGroups {
		// If the pattern is "*", return all resource groups
		if pattern == "*" {
			plugin.Logger(ctx).Debug("listQueryResourceGroupsForConnection", "connection_name", d.Connection.Name, "pattern", pattern, "matching all resource groups")
			return resourceGroups, nil
		}

		// Match against the pattern
		for _, validResourceGroup := range resourceGroups {
			if ok, _ := path.Match(pattern, validResourceGroup); ok {
				targetResourceGroups = append(targetResourceGroups, validResourceGroup)
			}
		}
	}

	// Remove duplicates
	targetResourceGroups = helpers.StringSliceDistinct(targetResourceGroups)

	plugin.Logger(ctx).Debug("listQueryResourceGroupsForConnection", "connection_name", d.Connection.Name, "targetResourceGroups", targetResourceGroups)

	return targetResourceGroups, nil
}

// Get all resource groups in the subscription by using the existing listResourceGroups function
func getAllResourceGroups(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	// Create resources client
	resourcesClient := resources.NewGroupsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	resourcesClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &resourcesClient, d.Connection)

	var resourceGroups []string

	// List resource groups directly
	result, err := resourcesClient.List(ctx, "", nil)
	if err != nil {
		plugin.Logger(ctx).Error("getAllResourceGroups", "connection_name", d.Connection.Name, "error", err)
		return nil, err
	}

	// Process resource group results
	for _, resourceGroup := range result.Values() {
		if resourceGroup.Name != nil {
			resourceGroups = append(resourceGroups, strings.ToLower(*resourceGroup.Name))
		}
	}

	// Paginate through results
	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("getAllResourceGroups", "connection_name", d.Connection.Name, "error", err)
			return nil, err
		}
		for _, resourceGroup := range result.Values() {
			if resourceGroup.Name != nil {
				resourceGroups = append(resourceGroups, strings.ToLower(*resourceGroup.Name))
			}
		}
	}

	return resourceGroups, nil
}

// Helper function to filter by resource group at the table level
func filterResourceGroupFromID(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	id := d.Value.(string)
	parts := strings.Split(id, "/")
	for i := 0; i < len(parts)-1; i++ {
		if strings.EqualFold(parts[i], "resourceGroups") || strings.EqualFold(parts[i], "resourcegroups") {
			return strings.ToLower(parts[i+1]), nil
		}
	}

	return nil, nil
}

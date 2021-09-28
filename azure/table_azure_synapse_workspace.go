package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/synapse/mgmt/2021-03-01/synapse"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureSynapseWorkspace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_synapse_workspace",
		Description: "Azure Synapse Workspace",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getSynapseWorkspace,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listSynapseWorkspace,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Fully qualified resource ID for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "adla_resource_id",
				Description: "The ADLA resource ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.AdlaResourceID"),
			},
			{
				Name:        "managed_resource_group_name",
				Description: "The managed resource group of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.ManagedResourceGroupName"),
			},
			{
				Name:        "managed_virtual_network",
				Description: "A managed virtual network for the workspace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.ManagedVirtualNetwork"),
			},
			{
				Name:        "public_network_access",
				Description: "Pubic network access to workspace.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.PublicNetworkAccess"),
			},
			{
				Name:        "sql_administrator_login",
				Description: "Login for workspace SQL active directory administrator.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.SQLAdministratorLogin"),
			},
			{
				Name:        "sql_administrator_login_password",
				Description: "The SQL administrator login password of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("WorkspaceProperties.SQLAdministratorLoginPassword"),
			},
			{
				Name:        "connectivity_endpoints",
				Description: "Connectivity endpoints of the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkspaceProperties.ConnectivityEndpoints"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the resource.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listSynapseWorkspaceDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "default_data_lake_storage",
				Description: "Workspace default data lake storage account details.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkspaceProperties.DefaultDataLakeStorage"),
			},
			{
				Name:        "encryption",
				Description: "The encryption details of the workspace.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractSynapseWorkspaceEncryption),
			},
			{
				Name:        "extra_properties",
				Description: "Workspace level configs and feature flags.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkspaceProperties.ExtraProperties"),
			},
			{
				Name:        "identity",
				Description: "The identity of the workspace.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "managed_virtual_network_settings",
				Description: "Managed virtual network settings of the workspace.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkspaceProperties.ManagedVirtualNetworkSettings"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "Private endpoint connections to the workspace.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractSynapseWorkspacePrivateEndpointConnections),
			},
			{
				Name:        "purview_configuration",
				Description: "Purview configuration of the workspace.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkspaceProperties.PurviewConfiguration"),
			},
			{
				Name:        "virtual_network_profile",
				Description: "Virtual network profile of the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkspaceProperties.VirtualNetworkProfile"),
			},
			{
				Name:        "workspace_repository_configuration",
				Description: "Git integration settings of the workspace.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkspaceProperties.WorkspaceRepositoryConfiguration"),
			},
			{
				Name:        "workspace_uid",
				Description: "The unique identifier of the workspace.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("WorkspaceProperties.WorkspaceUID"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location").Transform(toLower),
			},
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(idToSubscriptionID),
			},
		},
	}
}

type SynapseWorkspaceEncryption struct {
	DoubleEncryptionEnabled *bool
	CmkStatus               *string
	CmkKey                  interface{}
}

//// LIST FUNCTION

func listSynapseWorkspace(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := synapse.NewWorkspacesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listSynapseWorkspace", "list", err)
		return nil, err
	}

	for _, config := range result.Values() {
		d.StreamListItem(ctx, config)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listSynapseWorkspace", "list_paging", err)
			return nil, err
		}
		for _, config := range result.Values() {
			d.StreamListItem(ctx, config)
		}
	}
	
	return nil, err
}

//// HYDRATE FUNCTIONS

func getSynapseWorkspace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSynapseWorkspace")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Handle empty name or resourceGroup
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := synapse.NewWorkspacesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	config, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("getSynapseWorkspace", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if config.ID != nil {
		return config, nil
	}

	return nil, nil
}

func listSynapseWorkspaceDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAppConfigurationDiagnosticSettings")
	id := *h.Item.(synapse.Workspace).ID

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewDiagnosticSettingsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("listAppConfigurationDiagnosticSettings", "list", err)
		return nil, err
	}

	// If we return the API response directly, the output does not provide
	// the contents of DiagnosticSettings
	var diagnosticSettings []map[string]interface{}
	for _, i := range *op.Value {
		objectMap := make(map[string]interface{})
		if i.ID != nil {
			objectMap["id"] = i.ID
		}
		if i.Name != nil {
			objectMap["name"] = i.Name
		}
		if i.Type != nil {
			objectMap["type"] = i.Type
		}
		if i.DiagnosticSettings != nil {
			objectMap["properties"] = i.DiagnosticSettings
		}
		diagnosticSettings = append(diagnosticSettings, objectMap)
	}
	return diagnosticSettings, nil
}

//// TRANSFORM FUNCTIONS

// If we return the API response directly, the output will not provide the properties of PrivateEndpointConnections
func extractSynapseWorkspacePrivateEndpointConnections(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	workspace := d.HydrateItem.(synapse.Workspace)
	var properties []map[string]interface{}

	if workspace.WorkspaceProperties.PrivateEndpointConnections != nil {
		for _, i := range *workspace.WorkspaceProperties.PrivateEndpointConnections {
			objectMap := make(map[string]interface{})
			if i.ID != nil {
				objectMap["id"] = i.ID
			}
			if i.ID != nil {
				objectMap["name"] = i.Name
			}
			if i.ID != nil {
				objectMap["type"] = i.Type
			}
			if i.PrivateEndpointConnectionProperties != nil {
				if i.PrivateEndpointConnectionProperties.PrivateEndpoint != nil {
					objectMap["privateEndpointPropertyId"] = i.PrivateEndpointConnectionProperties.PrivateEndpoint.ID
				}
				if i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState != nil {
					if i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.ActionsRequired != nil {
						objectMap["privateLinkServiceConnectionStateActionsRequired"] = i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.ActionsRequired
					}
					if i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.Status != nil {
						objectMap["privateLinkServiceConnectionStateStatus"] = i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.Status
					}
					if i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.Description != nil {
						objectMap["privateLinkServiceConnectionStateDescription"] = i.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState.Description
					}
				}
				if i.PrivateEndpointConnectionProperties.ProvisioningState != nil {
					objectMap["provisioningState"] = i.PrivateEndpointConnectionProperties.ProvisioningState
				}
			}
			properties = append(properties, objectMap)
		}
	}

	return properties, nil
}

// If we return the API response directly, the output will not provide the properties of Encryption
func extractSynapseWorkspaceEncryption(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	workspace := d.HydrateItem.(synapse.Workspace)
	var properties SynapseWorkspaceEncryption

	if workspace.WorkspaceProperties.Encryption != nil {
		if workspace.WorkspaceProperties.Encryption.DoubleEncryptionEnabled != nil {
			properties.DoubleEncryptionEnabled = workspace.WorkspaceProperties.Encryption.DoubleEncryptionEnabled
		}
		if workspace.WorkspaceProperties.Encryption.Cmk != nil {
			if workspace.WorkspaceProperties.Encryption.Cmk.Status != nil {
				properties.CmkStatus = workspace.WorkspaceProperties.Encryption.Cmk.Status
			}
			if workspace.WorkspaceProperties.Encryption.Cmk.Key != nil {
				properties.CmkKey = workspace.WorkspaceProperties.Encryption.Cmk.Key
			}
		}
	}

	return properties, nil
}

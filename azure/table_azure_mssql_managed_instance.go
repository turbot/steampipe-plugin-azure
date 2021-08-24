package azure

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql"
)

//// TABLE DEFINITION

func tableAzureMSSQLManagedInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_mssql_managed_instance",
		Description: "Azure Microsoft SQL Managed Instance",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getMSSQLManagedInstance,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404", "InvalidApiVersionParameter"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listMSSQLManagedInstances,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the managed instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a managed instance uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type of the managed instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.State"),
			},
			{
				Name:        "administrator_login",
				Description: "Administrator username for the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.AdministratorLogin"),
			},
			{
				Name:        "administrator_login_password",
				Description: "Administrator password for the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.AdministratorLoginPassword"),
			},
			{
				Name:        "collation",
				Description: "Collation of the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.Collation"),
			},
			{
				Name:        "dns_zone",
				Description: "The Dns zone that the managed instance is in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.DNSZone"),
			},
			{
				Name:        "dns_zone_partner",
				Description: "The resource id of another managed instance whose DNS zone this managed instance will share after creation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.DNSZonePartner"),
			},
			{
				Name:        "fully_qualified_domain_name",
				Description: "The fully qualified domain name of the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.FullyQualifiedDomainName"),
			},
			{
				Name:        "instance_pool_id",
				Description: "The Id of the instance pool this managed server belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.InstancePoolID"),
			},
			{
				Name:        "license_type",
				Description: "The license type of the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.LicenseType"),
			},
			{
				Name:        "maintenance_configuration_id",
				Description: "Specifies maintenance configuration id to apply to this managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.MaintenanceConfigurationID"),
			},
			{
				Name:        "managed_instance_create_mode",
				Description: "Specifies the mode of database creation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.ManagedInstanceCreateMode"),
			},
			{
				Name:        "minimal_tls_version",
				Description: "Minimal TLS version of the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.MinimalTLSVersion"),
			},
			{
				Name:        "proxy_override",
				Description: "Connection type used for connecting to the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.ProxyOverride"),
			},
			{
				Name:        "public_data_endpoint_enabled",
				Description: "Whether or not the public data endpoint is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ManagedInstanceProperties.PublicDataEndpointEnabled"),
			},
			{
				Name:        "restore_point_in_time",
				Description: "Specifies the point in time of the source database that will be restored to create the new database.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ManagedInstanceProperties.RestorePointInTime").Transform(convertDateToTime),
			},
			{
				Name:        "source_managed_instance_id",
				Description: "The resource identifier of the source managed instance associated with create operation of this instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.SourceManagedInstanceID"),
			},
			{
				Name:        "storage_size_in_gb",
				Description: "The managed instance storage size in GB.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ManagedInstanceProperties.StorageSizeInGB"),
			},
			{
				Name:        "subnet_id",
				Description: "Subnet resource ID for the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.SubnetID"),
			},
			{
				Name:        "timezone_id",
				Description: "Id of the timezone.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.TimezoneID"),
			},
			{
				Name:        "v_cores",
				Description: "The number of vcores of the managed instance.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ManagedInstanceProperties.VCores"),
			},
			{
				Name:        "encryption_protectors",
				Description: "The managed instance encryption protectors.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listMSSQLManagedInstanceEncryptionProtectors,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "identity",
				Description: "The azure active directory identity of the managed instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "sku",
				Description: "Managed instance SKU.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vulnerability_assessments",
				Description: "The managed instance vulnerability assessments.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listMSSQLManagedInstanceVulnerabilityAssessments,
				Transform:   transform.FromValue(),
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

//// LIST FUNCTION

func listMSSQLManagedInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := sql.NewManagedInstancesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx, "")
	if err != nil {
		return nil, err
	}
	for _, managedInstance := range result.Values() {
		d.StreamListItem(ctx, managedInstance)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, managedInstance := range result.Values() {
			d.StreamListItem(ctx, managedInstance)
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getMSSQLManagedInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMSSQLManagedInstance")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, of no input provided
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := sql.NewManagedInstancesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

func listMSSQLManagedInstanceEncryptionProtectors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listMSSQLManagedInstanceEncryptionProtectors")

	managedInstance := h.Item.(sql.ManagedInstance)
	resourceGroup := strings.Split(string(*managedInstance.ID), "/")[4]
	managedInstanceName := *managedInstance.Name

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := sql.NewManagedInstanceEncryptionProtectorsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByInstance(ctx, resourceGroup, managedInstanceName)
	if err != nil {
		return nil, err
	}

	var managedInstanceEncryptionProtectors []map[string]interface{}

	for _, i := range op.Values() {
		managedInstanceEncryptionProtector := make(map[string]interface{})
		if i.ID != nil {
			managedInstanceEncryptionProtector["id"] = *i.ID
		}
		if i.Name != nil {
			managedInstanceEncryptionProtector["name"] = *i.Name
		}
		if i.Type != nil {
			managedInstanceEncryptionProtector["type"] = *i.Type
		}
		if i.Kind != nil {
			managedInstanceEncryptionProtector["kind"] = *i.Kind
		}
		if i.ManagedInstanceEncryptionProtectorProperties.AutoRotationEnabled != nil {
			managedInstanceEncryptionProtector["autoRotationEnabled"] = i.ManagedInstanceEncryptionProtectorProperties.AutoRotationEnabled
		}
		if i.ManagedInstanceEncryptionProtectorProperties.ServerKeyName != nil {
			managedInstanceEncryptionProtector["serverKeyName"] = i.ManagedInstanceEncryptionProtectorProperties.ServerKeyName
		}
		if len(i.ManagedInstanceEncryptionProtectorProperties.ServerKeyType) > 0 {
			managedInstanceEncryptionProtector["serverKeyType"] = i.ManagedInstanceEncryptionProtectorProperties.ServerKeyType
		}
		if i.ManagedInstanceEncryptionProtectorProperties.Thumbprint != nil {
			managedInstanceEncryptionProtector["thumbprint"] = i.ManagedInstanceEncryptionProtectorProperties.Thumbprint
		}
		if i.ManagedInstanceEncryptionProtectorProperties.URI != nil {
			managedInstanceEncryptionProtector["uri"] = i.ManagedInstanceEncryptionProtectorProperties.URI
		}

		managedInstanceEncryptionProtectors = append(managedInstanceEncryptionProtectors, managedInstanceEncryptionProtector)
	}

	for op.NotDone() {
		err = op.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, i := range op.Values() {
			managedInstanceEncryptionProtector := make(map[string]interface{})
			if i.ID != nil {
				managedInstanceEncryptionProtector["id"] = *i.ID
			}
			if i.Name != nil {
				managedInstanceEncryptionProtector["name"] = *i.Name
			}
			if i.Type != nil {
				managedInstanceEncryptionProtector["type"] = *i.Type
			}
			if i.Kind != nil {
				managedInstanceEncryptionProtector["kind"] = *i.Kind
			}
			if i.ManagedInstanceEncryptionProtectorProperties.AutoRotationEnabled != nil {
				managedInstanceEncryptionProtector["autoRotationEnabled"] = i.ManagedInstanceEncryptionProtectorProperties.AutoRotationEnabled
			}
			if i.ManagedInstanceEncryptionProtectorProperties.ServerKeyName != nil {
				managedInstanceEncryptionProtector["serverKeyName"] = i.ManagedInstanceEncryptionProtectorProperties.ServerKeyName
			}
			if len(i.ManagedInstanceEncryptionProtectorProperties.ServerKeyType) > 0 {
				managedInstanceEncryptionProtector["serverKeyType"] = i.ManagedInstanceEncryptionProtectorProperties.ServerKeyType
			}
			if i.ManagedInstanceEncryptionProtectorProperties.Thumbprint != nil {
				managedInstanceEncryptionProtector["thumbprint"] = i.ManagedInstanceEncryptionProtectorProperties.Thumbprint
			}
			if i.ManagedInstanceEncryptionProtectorProperties.URI != nil {
				managedInstanceEncryptionProtector["uri"] = i.ManagedInstanceEncryptionProtectorProperties.URI
			}

			managedInstanceEncryptionProtectors = append(managedInstanceEncryptionProtectors, managedInstanceEncryptionProtector)
		}
	}

	return managedInstanceEncryptionProtectors, nil
}

func listMSSQLManagedInstanceVulnerabilityAssessments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listMSSQLManagedInstanceVulnerabilityAssessments")

	managedInstance := h.Item.(sql.ManagedInstance)
	resourceGroup := strings.Split(string(*managedInstance.ID), "/")[4]
	managedInstanceName := *managedInstance.Name

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := sql.NewManagedInstanceVulnerabilityAssessmentsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByInstance(ctx, resourceGroup, managedInstanceName)
	if err != nil {
		return nil, err
	}

	var managedInstanceVulnerabilityAssessments []map[string]interface{}

	for _, i := range op.Values() {
		managedInstanceVulnerabilityAssessment := make(map[string]interface{})
		if i.ID != nil {
			managedInstanceVulnerabilityAssessment["id"] = *i.ID
		}
		if i.Name != nil {
			managedInstanceVulnerabilityAssessment["name"] = *i.Name
		}
		if i.Type != nil {
			managedInstanceVulnerabilityAssessment["type"] = *i.Type
		}
		if i.ManagedInstanceVulnerabilityAssessmentProperties.RecurringScans != nil {
			managedInstanceVulnerabilityAssessment["recurringScans"] = i.ManagedInstanceVulnerabilityAssessmentProperties.RecurringScans
		}
		if i.ManagedInstanceVulnerabilityAssessmentProperties.StorageAccountAccessKey != nil {
			managedInstanceVulnerabilityAssessment["storageAccountAccessKey"] = *i.ManagedInstanceVulnerabilityAssessmentProperties.StorageAccountAccessKey
		}
		if i.ManagedInstanceVulnerabilityAssessmentProperties.StorageContainerPath != nil {
			managedInstanceVulnerabilityAssessment["storageContainerPath"] = *i.ManagedInstanceVulnerabilityAssessmentProperties.StorageContainerPath
		}
		if i.ManagedInstanceVulnerabilityAssessmentProperties.StorageContainerSasKey != nil {
			managedInstanceVulnerabilityAssessment["storageContainerSasKey"] = *i.ManagedInstanceVulnerabilityAssessmentProperties.StorageContainerSasKey
		}

		managedInstanceVulnerabilityAssessments = append(managedInstanceVulnerabilityAssessments, managedInstanceVulnerabilityAssessment)
	}

	for op.NotDone() {
		err = op.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, i := range op.Values() {
			managedInstanceVulnerabilityAssessment := make(map[string]interface{})
			if i.ID != nil {
				managedInstanceVulnerabilityAssessment["id"] = *i.ID
			}
			if i.Name != nil {
				managedInstanceVulnerabilityAssessment["name"] = *i.Name
			}
			if i.Type != nil {
				managedInstanceVulnerabilityAssessment["type"] = *i.Type
			}
			if i.ManagedInstanceVulnerabilityAssessmentProperties.RecurringScans != nil {
				managedInstanceVulnerabilityAssessment["recurringScans"] = i.ManagedInstanceVulnerabilityAssessmentProperties.RecurringScans
			}
			if i.ManagedInstanceVulnerabilityAssessmentProperties.StorageAccountAccessKey != nil {
				managedInstanceVulnerabilityAssessment["storageAccountAccessKey"] = *i.ManagedInstanceVulnerabilityAssessmentProperties.StorageAccountAccessKey
			}
			if i.ManagedInstanceVulnerabilityAssessmentProperties.StorageContainerPath != nil {
				managedInstanceVulnerabilityAssessment["storageContainerPath"] = *i.ManagedInstanceVulnerabilityAssessmentProperties.StorageContainerPath
			}
			if i.ManagedInstanceVulnerabilityAssessmentProperties.StorageContainerSasKey != nil {
				managedInstanceVulnerabilityAssessment["storageContainerSasKey"] = *i.ManagedInstanceVulnerabilityAssessmentProperties.StorageContainerSasKey
			}

			managedInstanceVulnerabilityAssessments = append(managedInstanceVulnerabilityAssessments, managedInstanceVulnerabilityAssessment)
		}
	}

	return managedInstanceVulnerabilityAssessments, nil
}

package azure

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
)

//// TABLE DEFINITION

func tableAzureMSSQLManagedInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_mssql_managed_instance",
		Description: "Azure Microsoft SQL Managed Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getMSSQLManagedInstance,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404", "InvalidApiVersionParameter"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listMSSQLManagedInstances,
		},
		Columns: azureColumns([]*plugin.Column{
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
				Transform:   transform.FromField("Properties.State"),
			},
			{
				Name:        "administrator_login",
				Description: "Administrator username for the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.AdministratorLogin"),
			},
			{
				Name:        "administrator_login_password",
				Description: "Administrator password for the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.AdministratorLoginPassword"),
			},
			{
				Name:        "collation",
				Description: "Collation of the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Collation"),
			},
			{
				Name:        "dns_zone",
				Description: "The Dns zone that the managed instance is in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DNSZone"),
			},
			{
				Name:        "dns_zone_partner",
				Description: "The resource id of another managed instance whose DNS zone this managed instance will share after creation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DNSZonePartner"),
			},
			{
				Name:        "fully_qualified_domain_name",
				Description: "The fully qualified domain name of the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.FullyQualifiedDomainName"),
			},
			{
				Name:        "instance_pool_id",
				Description: "The Id of the instance pool this managed server belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.InstancePoolID"),
			},
			{
				Name:        "license_type",
				Description: "The license type of the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.LicenseType"),
			},
			{
				Name:        "maintenance_configuration_id",
				Description: "Specifies maintenance configuration id to apply to this managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.MaintenanceConfigurationID"),
			},
			{
				Name:        "managed_instance_create_mode",
				Description: "Specifies the mode of database creation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ManagedInstanceCreateMode"),
			},
			{
				Name:        "minimal_tls_version",
				Description: "Minimal TLS version of the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.MinimalTLSVersion"),
			},
			{
				Name:        "proxy_override",
				Description: "Connection type used for connecting to the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProxyOverride"),
			},
			{
				Name:        "public_data_endpoint_enabled",
				Description: "Whether or not the public data endpoint is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.PublicDataEndpointEnabled"),
			},
			{
				Name:        "restore_point_in_time",
				Description: "Specifies the point in time of the source database that will be restored to create the new database.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.RestorePointInTime"),
			},
			{
				Name:        "source_managed_instance_id",
				Description: "The resource identifier of the source managed instance associated with create operation of this instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.SourceManagedInstanceID"),
			},
			{
				Name:        "storage_size_in_gb",
				Description: "The managed instance storage size in GB.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.StorageSizeInGB"),
			},
			{
				Name:        "subnet_id",
				Description: "Subnet resource ID for the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.SubnetID"),
			},
			{
				Name:        "timezone_id",
				Description: "Id of the timezone.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.TimezoneID"),
			},
			{
				Name:        "v_cores",
				Description: "The number of vcores of the managed instance.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.VCores"),
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
				Name:        "security_alert_policies",
				Description: "The security alert policies of the managed instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listMSSQLManagedInstanceSecurityAlertPolicies,
				Transform:   transform.FromValue(),
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
		}),
	}
}

//// LIST FUNCTION

func listMSSQLManagedInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}

	client, err := armsql.NewManagedInstancesClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}

	pager := client.NewListPager(nil)
	for pager.More() {
		result, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, managedInstance := range result.Value {
			d.StreamListItem(ctx, *managedInstance)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getMSSQLManagedInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getMSSQLManagedInstance")

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// Return nil, of no input provided
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	client, err := armsql.NewManagedInstancesClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}

	op, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		plugin.Logger(ctx).Error("getMSSQLManagedInstance", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op.ManagedInstance, nil
	}

	return nil, nil
}

func listMSSQLManagedInstanceEncryptionProtectors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listMSSQLManagedInstanceEncryptionProtectors")

	managedInstance := h.Item.(armsql.ManagedInstance)
	resourceGroup := strings.Split(string(*managedInstance.ID), "/")[4]
	managedInstanceName := *managedInstance.Name

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	client, err := armsql.NewManagedInstanceEncryptionProtectorsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}

	pager := client.NewListByInstancePager(resourceGroup, managedInstanceName, nil)
	var managedInstanceEncryptionProtectors []*armsql.ManagedInstanceEncryptionProtector
	for pager.More() {
		result, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		managedInstanceEncryptionProtectors = append(managedInstanceEncryptionProtectors, result.Value...)
	}

	return managedInstanceEncryptionProtectors, nil
}

func listMSSQLManagedInstanceVulnerabilityAssessments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listMSSQLManagedInstanceVulnerabilityAssessments")

	managedInstance := h.Item.(armsql.ManagedInstance)
	resourceGroup := strings.Split(string(*managedInstance.ID), "/")[4]
	managedInstanceName := *managedInstance.Name

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	client, err := armsql.NewManagedInstanceVulnerabilityAssessmentsClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}

	pager := client.NewListByInstancePager(resourceGroup, managedInstanceName, nil)
	var managedInstanceVulnerabilityAssessments []*armsql.ManagedInstanceVulnerabilityAssessment
	for pager.More() {
		result, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		managedInstanceVulnerabilityAssessments = append(managedInstanceVulnerabilityAssessments, result.Value...)
	}

	return managedInstanceVulnerabilityAssessments, nil
}

func listMSSQLManagedInstanceSecurityAlertPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listMSSQLManagedInstanceSecurityAlertPolicies")

	managedInstance := h.Item.(armsql.ManagedInstance)
	resourceGroup := strings.Split(string(*managedInstance.ID), "/")[4]
	managedInstanceName := *managedInstance.Name

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}
	client, err := armsql.NewManagedServerSecurityAlertPoliciesClient(session.SubscriptionID, session.Cred, session.ClientOptions)
	if err != nil {
		return nil, err
	}

	pager := client.NewListByInstancePager(resourceGroup, managedInstanceName, nil)
	var managedInstanceSecurityAlertPolicies []*armsql.ManagedServerSecurityAlertPolicy
	for pager.More() {
		result, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		managedInstanceSecurityAlertPolicies = append(managedInstanceSecurityAlertPolicies, result.Value...)
	}

	return managedInstanceSecurityAlertPolicies, nil
}

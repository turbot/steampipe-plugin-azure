package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/services/databox/mgmt/2020-11-01/databox"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureDataBoxJob(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_databox_job",
		Description: "Azure Data Box Job",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getDataBoxJob,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "400"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listResourceGroups,
			Hydrate:       listDataBoxJobs,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The resource identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "Type of the object.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sku_name",
				Description: "The sku name. Possible values include: 'DataBox', 'DataBoxDisk', 'DataBoxHeavy'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "sku_display_name",
				Description: "The display name of the sku.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.DisplayName"),
			},
			{
				Name:        "sku_Family",
				Description: "The sku family.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Family"),
			},
			{
				Name:        "transfer_type",
				Description: "Type of the data transfer. Possible values include: 'ImportToAzure', 'ExportFromAzure'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobProperties.TransferType"),
			},
			{
				Name:        "is_cancellable",
				Description: "Describes whether the job is cancellable or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("JobProperties.IsCancellable"),
			},
			{
				Name:        "is_deletable",
				Description: "Describes whether the job is deletable or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("JobProperties.IsDeletable"),
			},
			{
				Name:        "is_shipping_address_editable",
				Description: "Describes whether the shipping address is editable or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("JobProperties.IsShippingAddressEditable"),
			},
			{
				Name:        "is_prepare_to_ship_enabled",
				Description: "Is Prepare To Ship Enabled on this job.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("JobProperties.IsPrepareToShipEnabled"),
			},
			{
				Name:        "location",
				Description: "The location of the resource. This will be one of the supported and registered Azure Regions (e.g. West US, East US, Southeast Asia, etc.). The region of a resource cannot be changed once it is created, but if an identical region is specified on update the request will succeed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "Name of the stage which is in progress. Possible values include: 'StageNameDeviceOrdered', 'StageNameDevicePrepared', 'StageNameDispatched', 'StageNameDelivered', 'StageNamePickedUp', 'StageNameAtAzureDC', 'StageNameDataCopy', 'StageNameCompleted', 'StageNameCompletedWithErrors', 'StageNameCancelled', 'StageNameFailedIssueReportedAtCustomer', 'StageNameFailedIssueDetectedAtAzureDC', 'StageNameAborted', 'StageNameCompletedWithWarnings', 'StageNameReadyToDispatchFromAzureDC', 'StageNameReadyToReceiveAtAzureDC'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobProperties.Status"),
			},
			{
				Name:        "start_time",
				Description: "Time at which the job was started in UTC ISO 8601 format.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("JobProperties.StartTime.Time"),
			},
			{
				Name:        "cancellation_reason",
				Description: "Reason for cancellation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobProperties.CancellationReason"),
			},
			{
				Name:        "delivery_type",
				Description: "Delivery type of Job. Possible values include: 'NonScheduled', 'Scheduled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobProperties.DeliveryType"),
			},
			{
				Name:        "is_cancellable_without_fee",
				Description: "Flag to indicate cancellation of scheduled job.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("JobProperties.IsCancellableWithoutFee"),
			},
			{
				Name:        "error",
				Description: "Top level error for the job.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("JobProperties.Error"),
			},
			{
				Name:        "details",
				Description: "Details of a job run. This field will only be sent for expand details filter.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("JobProperties.Details"),
			},
			{
				Name:        "delivery_info",
				Description: "Delivery Info of Job.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("JobProperties.DeliveryInfo"),
			},
			{
				Name:        "identity",
				Description: "Msi identity of the resource.",
				Type:        proto.ColumnType_JSON,
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

func listDataBoxJobs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_databox_job.listDataBoxJobs", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	deviceClient := databox.NewJobsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	deviceClient.Authorizer = session.Authorizer

	resourceGroupName := h.Item.(resources.Group).Name

	result, err := deviceClient.ListByResourceGroup(ctx, *resourceGroupName, "")
	if err != nil {
		plugin.Logger(ctx).Error("azure_databox_job.listDataBoxJobs", "api_error", err)
		return nil, err
	}
	for _, device := range result.Values() {
		d.StreamListItem(ctx, device)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_databox_job.listDataBoxJobs", "pagination_error", err)
			return nil, err
		}
		for _, job := range result.Values() {
			d.StreamListItem(ctx, job)
		}

	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getDataBoxJob(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	// Return nil, if no input provide
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_databox_job.getDataBoxJob", "session_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	deviceClient := databox.NewJobsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	deviceClient.Authorizer = session.Authorizer

	op, err := deviceClient.Get(ctx, resourceGroup, name, "")
	if err != nil {
		plugin.Logger(ctx).Error("azure_databox_job.getDataBoxJob", "api_error", err)
		return nil, err
	}

	return op, nil
}

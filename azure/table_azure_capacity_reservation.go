package azure

import (
	"context"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/reservations/armreservations"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureCapacityReservation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_capacity_reservation",
		Description: "Azure Capacity Reservation",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"reservation_order_id", "reservation_id"}),
			Hydrate:    getAzureCapacityReservation,
			Tags: map[string]string{
				"service": "Microsoft.Capacity",
				"action":  "reservations/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ReservationOrderIdInvalid", "ReservationOrderNotFound", "ReservationNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureCapacityReservations,
			Tags: map[string]string{
				"service": "Microsoft.Capacity",
				"action":  "reservations/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The fully qualified ID for the reservation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "name",
				Description: "The name of the reservation in format {reservationOrderId}/{reservationId}.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The resource type of the reservation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "reservation_order_id",
				Description: "The ID of the reservation order that contains this reservation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(extractReservationOrderID),
			},
			{
				Name:        "reservation_id",
				Description: "The ID of the individual reservation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name").Transform(extractReservationID),
			},
			{
				Name:        "sku_name",
				Description: "The name of the SKU for the reservation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SKU.Name"),
			},
			{
				Name:        "kind",
				Description: "The kind of the reservation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "The etag of the reservation.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "display_name",
				Description: "Friendly name for user to easily identify the reservation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DisplayName"),
			},
			{
				Name:        "provisioning_state",
				Description: "Current state of the reservation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "provisioning_sub_state",
				Description: "The provisioning sub-state of the reservation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningSubState"),
			},
			{
				Name:        "reserved_resource_type",
				Description: "The type of the resource that is being reserved.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ReservedResourceType"),
			},
			{
				Name:        "quantity",
				Description: "The number of instances reserved in the reservation.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.Quantity"),
			},
			{
				Name:        "term",
				Description: "The reservation term (P1Y, P3Y, P5Y).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Term"),
			},
			{
				Name:        "billing_plan",
				Description: "The billing plan for the reservation (Upfront or Monthly).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.BillingPlan"),
			},
			{
				Name:        "billing_scope_id",
				Description: "Billing scope of the reservation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.BillingScopeID"),
			},
			{
				Name:        "applied_scope_type",
				Description: "The scope type of the applied reservation (Single, Shared, ManagementGroup).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.AppliedScopeType"),
			},
			{
				Name:        "applied_scopes",
				Description: "List of the subscriptions that the benefit will be applied to.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.AppliedScopes"),
			},
			{
				Name:        "instance_flexibility",
				Description: "Allows reservation discount to be applied across skus within the same Autofit group (On or Off).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.InstanceFlexibility"),
			},
			{
				Name:        "benefit_start_time",
				Description: "The time when the reservation benefit starts.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.BenefitStartTime").Transform(reservationTimeToTimestamp),
			},
			{
				Name:        "purchase_date",
				Description: "The date when the reservation was purchased.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.PurchaseDate").Transform(reservationTimeToTimestamp),
			},
			{
				Name:        "expiry_date",
				Description: "The date when the reservation expires.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.ExpiryDate").Transform(reservationTimeToTimestamp),
			},
			{
				Name:        "last_updated_date_time",
				Description: "The date and time the reservation was last updated.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.LastUpdatedDateTime").Transform(reservationTimeToTimestamp),
			},
			{
				Name:        "renew",
				Description: "Setting this to true will automatically purchase a new reservation on the expiry date.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.Renew"),
			},
			{
				Name:        "renew_source",
				Description: "The reservation ID that is the source of the renewal.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.RenewSource"),
			},
			{
				Name:        "renew_destination",
				Description: "The reservation ID that replaces this reservation when it expires.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.RenewDestination"),
			},
			{
				Name:        "archived",
				Description: "Indicates if the reservation is archived.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.Archived"),
			},
			{
				Name:        "sku_description",
				Description: "Description of the SKU in english.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.SKUDescription"),
			},
			{
				Name:        "utilization",
				Description: "Reservation utilization.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Utilization"),
			},
			{
				Name:        "extended_status_info",
				Description: "Additional information about the status of the reservation.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.ExtendedStatusInfo"),
			},
			{
				Name:        "merge_properties",
				Description: "Properties of the reservation for merge operations.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.MergeProperties"),
			},
			{
				Name:        "split_properties",
				Description: "Properties of the reservation for split operations.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.SplitProperties"),
			},
			{
				Name:        "system_data",
				Description: "The system metadata relating to this reservation.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DisplayName"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},
			// Note: no resource_group (reservations are not in a resource group)
			// Note: no tags (not supported by this resource type)
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location").Transform(toLower),
			},
		}),
	}
}

//// LIST FUNCTION

func listAzureCapacityReservations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_capacity_reservation.listAzureCapacityReservations", "session_error", err)
		return nil, err
	}

	client, err := armreservations.NewReservationClient(
		session.Cred,
		session.ClientOptions,
	)
	if err != nil {
		plugin.Logger(ctx).Error("azure_capacity_reservation.listAzureCapacityReservations", "client_error", err)
		return nil, err
	}

	pager := client.NewListAllPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			// Azure API sometimes returns "value": "" (string) instead of "value": []
			// for tenants with no reservations. Treat as empty result.
			if strings.Contains(err.Error(), "cannot unmarshal string into Go value of type []*armreservations.ReservationResponse") {
				plugin.Logger(ctx).Warn("azure_capacity_reservation.listAzureCapacityReservations",
					"msg", "Azure API returned malformed response for empty reservation list, treating as empty")
				return nil, nil
			}
			plugin.Logger(ctx).Error("azure_capacity_reservation.listAzureCapacityReservations", "list_error", err)
			return nil, err
		}
		for _, item := range page.Value {
			d.StreamListItem(ctx, item)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// GET FUNCTION

func getAzureCapacityReservation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	reservationOrderID := d.EqualsQualString("reservation_order_id")
	reservationID := d.EqualsQualString("reservation_id")

	if reservationOrderID == "" || reservationID == "" {
		return nil, nil
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_capacity_reservation.getAzureCapacityReservation", "session_error", err)
		return nil, err
	}

	client, err := armreservations.NewReservationClient(
		session.Cred,
		session.ClientOptions,
	)
	if err != nil {
		plugin.Logger(ctx).Error("azure_capacity_reservation.getAzureCapacityReservation", "client_error", err)
		return nil, err
	}

	// Note: Get signature is Get(ctx, reservationID, reservationOrderID, options)
	result, err := client.Get(ctx, reservationID, reservationOrderID, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_capacity_reservation.getAzureCapacityReservation", "get_error", err)
		return nil, err
	}

	return &result.ReservationResponse, nil
}

//// TRANSFORM FUNCTIONS

// extractReservationOrderID extracts the order ID from a reservation Name field.
// The Name field format is "{reservationOrderId}/{reservationId}".
func extractReservationOrderID(_ context.Context, d *transform.TransformData) (interface{}, error) {
	name, ok := d.Value.(*string)
	if !ok || name == nil {
		return nil, nil
	}
	parts := strings.SplitN(*name, "/", 2)
	if len(parts) < 1 {
		return nil, nil
	}
	return parts[0], nil
}

// extractReservationID extracts the reservation ID from a reservation Name field.
// The Name field format is "{reservationOrderId}/{reservationId}".
func extractReservationID(_ context.Context, d *transform.TransformData) (interface{}, error) {
	name, ok := d.Value.(*string)
	if !ok || name == nil {
		return nil, nil
	}
	parts := strings.SplitN(*name, "/", 2)
	if len(parts) < 2 {
		return nil, nil
	}
	return parts[1], nil
}

// reservationTimeToTimestamp converts a *time.Time to RFC3339 string for Steampipe.
func reservationTimeToTimestamp(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}
	t, ok := d.Value.(*time.Time)
	if !ok || t == nil {
		return nil, nil
	}
	return t.UTC().Format(time.RFC3339), nil
}

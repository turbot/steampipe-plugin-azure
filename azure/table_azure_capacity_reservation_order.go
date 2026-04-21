package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/reservations/armreservations"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureCapacityReservationOrder(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_capacity_reservation_order",
		Description: "Azure Capacity Reservation Order",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getAzureCapacityReservationOrder,
			Tags: map[string]string{
				"service": "Microsoft.Capacity",
				"action":  "reservationOrders/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ReservationOrderIdInvalid", "ReservationOrderNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureCapacityReservationOrders,
			Tags: map[string]string{
				"service": "Microsoft.Capacity",
				"action":  "reservationOrders/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The fully qualified ID for the reservation order.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "name",
				Description: "The name of the reservation order (UUID).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The resource type of the reservation order.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "The etag of the reservation order.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "display_name",
				Description: "Friendly name for the reservation order.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.DisplayName"),
			},
			{
				Name:        "provisioning_state",
				Description: "Current provisioning state of the reservation order.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.ProvisioningState"),
			},
			{
				Name:        "billing_plan",
				Description: "The billing plan for the reservation order (Upfront or Monthly).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.BillingPlan"),
			},
			{
				Name:        "term",
				Description: "The reservation term (P1Y, P3Y, P5Y).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Term"),
			},
			{
				Name:        "original_quantity",
				Description: "Total quantity of the SKUs purchased in the reservation order.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.OriginalQuantity"),
			},
			{
				Name:        "benefit_start_time",
				Description: "The time when the reservation order benefit starts.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.BenefitStartTime").Transform(reservationTimeToTimestamp),
			},
			{
				Name:        "created_date_time",
				Description: "The date and time the reservation order was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.CreatedDateTime").Transform(reservationTimeToTimestamp),
			},
			{
				Name:        "expiry_date",
				Description: "The date when the reservation order expires.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.ExpiryDate").Transform(reservationTimeToTimestamp),
			},
			{
				Name:        "request_date_time",
				Description: "The date and time the reservation order was requested.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Properties.RequestDateTime").Transform(reservationTimeToTimestamp),
			},
			{
				Name:        "plan_information",
				Description: "Information describing the type of billing plan for this reservation order.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.PlanInformation"),
			},
			{
				Name:        "reservations",
				Description: "List of reservations in this reservation order.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.Reservations"),
			},
			{
				Name:        "system_data",
				Description: "The system metadata relating to this reservation order.",
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
			// Note: no resource_group, no tags, no region — reservation orders are tenant-level resources
		}),
	}
}

//// LIST FUNCTION

func listAzureCapacityReservationOrders(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_capacity_reservation_order.listAzureCapacityReservationOrders", "session_error", err)
		return nil, err
	}

	client, err := armreservations.NewReservationOrderClient(
		session.Cred,
		session.ClientOptions,
	)
	if err != nil {
		plugin.Logger(ctx).Error("azure_capacity_reservation_order.listAzureCapacityReservationOrders", "client_error", err)
		return nil, err
	}

	pager := client.NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			// Azure API sometimes returns "value": "" (string) instead of "value": []
			// for tenants with no reservation orders. Treat as empty result.
			if strings.Contains(err.Error(), "cannot unmarshal string into Go value of type []*armreservations.ReservationOrderResponse") {
				plugin.Logger(ctx).Warn("azure_capacity_reservation_order.listAzureCapacityReservationOrders",
					"msg", "Azure API returned malformed response for empty reservation order list, treating as empty")
				return nil, nil
			}
			plugin.Logger(ctx).Error("azure_capacity_reservation_order.listAzureCapacityReservationOrders", "list_error", err)
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

func getAzureCapacityReservationOrder(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	reservationOrderID := d.EqualsQualString("name")

	if reservationOrderID == "" {
		return nil, nil
	}

	session, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azure_capacity_reservation_order.getAzureCapacityReservationOrder", "session_error", err)
		return nil, err
	}

	client, err := armreservations.NewReservationOrderClient(
		session.Cred,
		session.ClientOptions,
	)
	if err != nil {
		plugin.Logger(ctx).Error("azure_capacity_reservation_order.getAzureCapacityReservationOrder", "client_error", err)
		return nil, err
	}

	// Expand planInformation to get full billing schedule details
	expand := "planInformation"
	result, err := client.Get(ctx, reservationOrderID, &armreservations.ReservationOrderClientGetOptions{
		Expand: &expand,
	})
	if err != nil {
		plugin.Logger(ctx).Error("azure_capacity_reservation_order.getAzureCapacityReservationOrder", "get_error", err)
		return nil, err
	}

	return &result.ReservationOrderResponse, nil
}

package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-10-01/consumption"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableAzureReservationRecommendation(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_reservation_recommendation",
		Description: "Azure Reservation Recommendation",
		List: &plugin.ListConfig{
			Hydrate: listReservedInstanceRecomendations,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The ID that uniquely identifies an event.",
			},
			{
				Name:        "id",
				Description: "The full qualified ARM ID of an event.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "kind",
				Description: "Specifies the kind of reservation recommendation.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "The etag for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "Resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sku",
				Description: "Resource sku.",
				Type:        proto.ColumnType_STRING,
			},

			// JSON fields
			{
				Name:        "legacy_recommendation_properties",
				Description: "The legacy recommendation properties.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "modern_recommendation_properties",
				Description: "The legacy recommendation properties.",
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
		}),
	}
}

type RecomendationInfo struct {
	LegacyRecommendationProperties *consumption.LegacyReservationRecommendationProperties
	ModernRecommendationProperties *consumption.ModernReservationRecommendationProperties
	ID                             *string
	Name                           *string
	Type                           *string
	Etag                           *string
	Tags                           map[string]*string
	Location                       *string
	Sku                            *string
	Kind                           consumption.KindBasicReservationRecommendation
}

//// LIST FUNCTION

func listReservedInstanceRecomendations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	reservedInstanceClient := consumption.NewReservationRecommendationsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	reservedInstanceClient.Authorizer = session.Authorizer
	result, err := reservedInstanceClient.List(ctx, "subscriptions/"+subscriptionID, "")
	if err != nil {
		return nil, err
	}
	for _, recomendation := range result.Values() {
		for _, r := range getReservationRecomendationProperties(recomendation) {
			d.StreamListItem(ctx, r)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, recomendation := range result.Values() {
			for _, r := range getReservationRecomendationProperties(recomendation) {
				d.StreamListItem(ctx, r)

				// Check if context has been cancelled or if the limit has been hit (if specified)
				// if there is a limit, it will return the number of rows required to reach this limit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, err
}

//// EXTRACT PROPERTIES

func getReservationRecomendationProperties(data consumption.BasicReservationRecommendation) []*RecomendationInfo {
	var results []*RecomendationInfo
	lInfo, isLegacy := data.AsLegacyReservationRecommendation()
	mInfo, isModern := data.AsModernReservationRecommendation()
	info, is := data.AsReservationRecommendation()
	if is {
		result := &RecomendationInfo{}
		result.Etag = info.Etag
		result.ID = info.ID
		result.Kind = info.Kind
		result.Location = info.Location
		result.Name = info.Name
		result.Sku = info.Sku
		result.Tags = info.Tags
		result.Type = info.Type
		results = append(results, result)
	}
	if isModern {
		result := &RecomendationInfo{}
		result.Etag = mInfo.Etag
		result.ID = mInfo.ID
		result.Kind = mInfo.Kind
		result.Location = mInfo.Location
		result.Name = mInfo.Name
		result.Sku = mInfo.Sku
		result.Tags = mInfo.Tags
		result.Type = mInfo.Type
		result.ModernRecommendationProperties = mInfo.ModernReservationRecommendationProperties
		results = append(results, result)
	}

	if isLegacy {
		result := &RecomendationInfo{}
		result.Etag = lInfo.Etag
		result.ID = lInfo.ID
		result.Kind = lInfo.Kind
		result.Location = lInfo.Location
		result.Name = lInfo.Name
		result.Sku = lInfo.Sku
		result.Tags = lInfo.Tags
		result.Type = lInfo.Type
		result.LegacyRecommendationProperties = lInfo.LegacyReservationRecommendationProperties
		results = append(results, result)
	}

	return results
}

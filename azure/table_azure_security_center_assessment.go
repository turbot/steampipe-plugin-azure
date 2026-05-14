package azure

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/security/mgmt/security"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

//// TABLE DEFINITION

func tableAzureSecurityCenterAssessment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_security_center_assessment",
		Description: "Azure Security Center Assessment",
		List: &plugin.ListConfig{
			Hydrate: listSecurityCenterAssessments,
			Tags: map[string]string{
				"service": "Microsoft.Security",
				"action":  "assessments/read",
			},
			KeyColumns: plugin.KeyColumnSlice{
				&plugin.KeyColumn{
					Name:      "expand",
					Operators: []string{"="},
					Require:   plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getSecurityCenterAssessment,
			Tags: map[string]string{
				"service": "Microsoft.Security",
				"action":  "assessments/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The resource id.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "User friendly display name of the assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssessmentPropertiesResponse.DisplayName"),
			},
			{
				Name:        "status_cause",
				Description: "Status cause of the assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssessmentPropertiesResponse.Status.Cause"),
			},
			{
				Name:        "status_code",
				Description: "Status code of the assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssessmentPropertiesResponse.Status.Code"),
			},
			{
				Name:        "status_description",
				Description: "Status description of the assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssessmentPropertiesResponse.Status.Description"),
			},
			{
				Name:        "status_change_date",
				Description: "Status change date of the assessment.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("AssessmentPropertiesResponse.Status.StatusChangeDate"),
			},
			{
				Name:        "status_first_evaluation_date",
				Description: "The timestamp of the first evaluation of the assessment.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("AssessmentPropertiesResponse.Status.FirstEvaluationDate"),
			},
			{
				Name:        "resource_details",
				Description: "Details of the resource that was assessed.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AssessmentPropertiesResponse.ResourceDetails"),
			},
			{
				Name:        "additional_data",
				Description: "Additional data regarding the assessment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AssessmentPropertiesResponse.AdditionalData"),
			},
			{
				Name:        "expand",
				Description: "May be used to expand the 'links' or 'metadata' of the assessment. By default, these fields are not included when listing or getting assessments.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("expand"),
			},
			{
				Name:        "links",
				Description: "Links relevant to the assessment. By default this is not populated, unless it's specified in expand.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AssessmentPropertiesResponse.Links"),
			},
			{
				Name:        "metadata",
				Description: "Describes properties of an assessment metadata. By default this is not populated, unless it's specified in expand.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AssessmentPropertiesResponse.Metadata"),
			},
			{
				Name:        "resource_name",
				Description: "Name of the resource that was assessed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssessmentPropertiesResponse.ResourceDetails").Transform(extractAssessmentResourceName),
			},
			{
				Name:        "resource_id",
				Description: "The full ID of the assessed resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractAssessedResourceIDFromAssessmentID),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssessmentPropertiesResponse.DisplayName"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},

			// Azure standard columns
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
		}),
	}
}

//// LOCAL TYPES

// The SDK's BasicResourceDetails unmarshaller switches on a lowercase "source" key, while the API
// returns PascalCase keys (Source, ResourceName, ...), so everything but Source gets dropped.
// To avoid that, responses are decoded into these local types, which keep resourceDetails as a map.
type securityCenterAssessment struct {
	ID                           *string `json:"id"`
	Name                         *string `json:"name"`
	Type                         *string `json:"type"`
	AssessmentPropertiesResponse struct {
		DisplayName *string `json:"displayName"`
		Status      *struct {
			Code                string     `json:"code"`
			Cause               string     `json:"cause"`
			Description         string     `json:"description"`
			StatusChangeDate    *time.Time `json:"statusChangeDate"`
			FirstEvaluationDate *time.Time `json:"firstEvaluationDate"`
		} `json:"status"`
		ResourceDetails map[string]interface{} `json:"resourceDetails"`
		AdditionalData  map[string]interface{} `json:"additionalData"`
		Links           map[string]interface{} `json:"links"`
		Metadata        map[string]interface{} `json:"metadata"`
	} `json:"properties"`
}

type securityCenterAssessmentList struct {
	Value    []securityCenterAssessment `json:"value"`
	NextLink *string                    `json:"nextLink"`
}

//// LIST FUNCTION

func listSecurityCenterAssessments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		logger.Error("azure_security_center_assessment.listSecurityCenterAssessments", "connection_error", err)
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	assessmentClient := security.NewAssessmentsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	assessmentClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &assessmentClient, d.Connection)

	req, err := assessmentClient.ListPreparer(ctx, "subscriptions/"+subscriptionID)
	if err != nil {
		logger.Error("azure_security_center_assessment.listSecurityCenterAssessments", "preparer_error", err)
		return nil, err
	}
	if expand := d.EqualsQualString("expand"); expand != "" {
		q := req.URL.Query()
		q.Set("$expand", expand)
		req.URL.RawQuery = q.Encode()
	}

	for {
		var page securityCenterAssessmentList
		if err := sendSecurityCenterAssessmentRequest(&assessmentClient, req, &page); err != nil {
			logger.Error("azure_security_center_assessment.listSecurityCenterAssessments", "query_error", err)
			return nil, err
		}

		for _, assessment := range page.Value {
			d.StreamListItem(ctx, assessment)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if page.NextLink == nil || *page.NextLink == "" {
			break
		}

		// Wait for rate limiting
		d.WaitForListRateLimit(ctx)

		req, err = http.NewRequestWithContext(ctx, http.MethodGet, *page.NextLink, nil)
		if err != nil {
			logger.Error("azure_security_center_assessment.listSecurityCenterAssessments", "next_page_error", err)
			return nil, err
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSecurityCenterAssessment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		logger.Error("azure_security_center_assessment.getSecurityCenterAssessment", "connection_error", err)
		return nil, err
	}

	assessmentID := d.EqualsQualString("id")
	if assessmentID == "" {
		return nil, nil
	}
	name := getLastPathElement(assessmentID)

	resourceID := assessmentID
	securityProviderIndex := strings.Index(strings.ToLower(assessmentID), "/providers/microsoft.security/assessments/")
	if securityProviderIndex > 0 {
		resourceID = assessmentID[:securityProviderIndex]
	}

	assessmentClient := security.NewAssessmentsClientWithBaseURI(session.ResourceManagerEndpoint, session.SubscriptionID)
	assessmentClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &assessmentClient, d.Connection)

	req, err := assessmentClient.GetPreparer(ctx, resourceID, name, getSecurityCenterAssessmentExpand(d))
	if err != nil {
		logger.Error("azure_security_center_assessment.getSecurityCenterAssessment", "preparer_error", err)
		return nil, err
	}

	var assessment securityCenterAssessment
	if err := sendSecurityCenterAssessmentRequest(&assessmentClient, req, &assessment); err != nil {
		logger.Error("azure_security_center_assessment.getSecurityCenterAssessment", "query_error", err)
		return nil, err
	}

	return assessment, nil
}

// sendSecurityCenterAssessmentRequest sends a prepared request through the client's sender
// (which handles retries) and decodes the JSON response body into result.
func sendSecurityCenterAssessmentRequest(client *security.AssessmentsClient, req *http.Request, result interface{}) error {
	resp, err := client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return err
	}
	return autorest.Respond(resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(result),
		autorest.ByClosing())
}

// getSecurityCenterAssessmentExpand builds the ExpandEnum value to send to the Security Center
// API from the "expand" qual, e.g. expand = 'links' or expand = 'links,metadata'.
func getSecurityCenterAssessmentExpand(d *plugin.QueryData) security.ExpandEnum {
	// In the get path "expand" isn't a key column, so read it from the raw quals.
	if q := d.Quals["expand"]; q != nil {
		for _, qual := range q.Quals {
			if qual.Operator == "=" && qual.Value != nil {
				return security.ExpandEnum(qual.Value.GetStringValue())
			}
		}
	}
	return security.ExpandEnum(d.EqualsQualString("expand"))
}

//// TRANSFORM FUNCTIONS

func extractAssessmentResourceName(_ context.Context, d *transform.TransformData) (interface{}, error) {
	resourceDetails, ok := d.Value.(map[string]interface{})
	if !ok {
		return nil, nil
	}
	for _, key := range []string{"ResourceName", "MachineName"} {
		if v, ok := resourceDetails[key].(string); ok && v != "" {
			return v, nil
		}
	}
	return nil, nil
}

func extractAssessedResourceIDFromAssessmentID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	assessmentID := types.SafeString(d.Value)
	if assessmentID == "" {
		return nil, nil
	}

	securityProviderIndex := strings.Index(strings.ToLower(assessmentID), "/providers/microsoft.security/assessments/")
	if securityProviderIndex > 0 {
		return assessmentID[:securityProviderIndex], nil
	}

	return assessmentID, nil
}

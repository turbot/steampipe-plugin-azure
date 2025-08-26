package azure

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/queueerror"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type queueInfo struct {
	Name          string
	Account       string
	ResourceGroup string
	Location      string
	Metadata      map[string]*string
	ID            string
	Type          string
}

//// TABLE DEFINITION

func tableAzureStorageQueue(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_storage_queue",
		Description: "Azure Storage Queue",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "storage_account_name", "resource_group"}),
			Hydrate:    getStorageQueue,
			Tags: map[string]string{
				"service": "Microsoft.Storage",
				"action":  "storageAccounts/queueServices/queues/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "QueueNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listStorageAccounts,
			Hydrate:       listStorageQueues,
			Tags: map[string]string{
				"service": "Microsoft.Storage",
				"action":  "storageAccounts/queueServices/queues/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The friendly name that identifies the queue."},
			{Name: "id", Description: "Resource ID of the queue.", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID")},
			{Name: "storage_account_name", Description: "The storage account containing the queue.", Type: proto.ColumnType_STRING, Transform: transform.FromField("Account")},
			{Name: "type", Description: "Type of the resource.", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type")},
			{Name: "metadata", Description: "Queue metadata key-value pairs.", Type: proto.ColumnType_JSON, Transform: transform.FromField("Metadata")},
			// Standard columns
			{Name: "title", Description: ColumnDescriptionTitle, Type: proto.ColumnType_STRING, Transform: transform.FromField("Name")},
			{Name: "akas", Description: ColumnDescriptionAkas, Type: proto.ColumnType_JSON, Transform: transform.From(buildQueueAkas)},
			// Azure standard
			{Name: "region", Description: ColumnDescriptionRegion, Type: proto.ColumnType_STRING, Transform: transform.FromField("Location").Transform(toLower)},
			{Name: "resource_group", Description: ColumnDescriptionResourceGroup, Type: proto.ColumnType_STRING, Transform: transform.FromField("ResourceGroup").Transform(toLower)},
		}),
	}
}

//// LIST FUNCTION

func listStorageQueues(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	acct := h.Item.(*storageAccountInfo)
	// Queue service not available for premium FileStorage or BlockBlobStorage accounts
	if acct.Account.Kind == "FileStorage" || acct.Account.Kind == "BlockBlobStorage" {
		return nil, nil
	}

	mSession, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	aadSession, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}

	allowShared := true
	if acct.Account.AccountProperties != nil && acct.Account.AccountProperties.AllowSharedKeyAccess != nil {
		allowShared = *acct.Account.AccountProperties.AllowSharedKeyAccess
	}

	config := GetConfig(d.Connection)
	authMode := "auto"
	if config.DataPlaneAuthMode != nil && *config.DataPlaneAuthMode != "" {
		authMode = strings.ToLower(*config.DataPlaneAuthMode)
	}

	qClient, err := buildQueueServiceClient(ctx, d, aadSession.Cred, authMode, *acct.Name, *acct.ResourceGroup, mSession.SubscriptionID, mSession.StorageEndpointSuffix, allowShared)
	if err != nil {
		return nil, err
	}

	pager := qClient.NewListQueuesPager(nil)
	for pager.More() {
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
		page, err := pager.NextPage(ctx)
		if err != nil {
			// Translate common auth errors and wrap properly
			if strings.Contains(err.Error(), "KeyBasedAuthenticationNotPermitted") {
				return nil, fmt.Errorf("shared key disabled; retry with auth_mode=aad: %w", err)
			}
			return nil, err
		}
		for _, q := range page.Queues {
			name := ""
			if q.Name != nil {
				name = *q.Name
			}
			meta := map[string]*string{}
			if q.Metadata != nil {
				meta = q.Metadata
			}
			id := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/queueServices/default/queues/%s", mSession.SubscriptionID, *acct.ResourceGroup, *acct.Name, name)
			qi := &queueInfo{Name: name, Account: *acct.Name, ResourceGroup: *acct.ResourceGroup, Location: *acct.Account.Location, Metadata: meta, ID: id, Type: "Microsoft.Storage/storageAccounts/queueServices/queues"}
			d.StreamListItem(ctx, qi)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getStorageQueue(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	accountName := d.EqualsQuals["storage_account_name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()
	name := d.EqualsQuals["name"].GetStringValue()
	if accountName == "" || resourceGroup == "" || name == "" {
		return nil, nil
	}

	mSession, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	aadSession, err := GetNewSessionUpdated(ctx, d)
	if err != nil {
		return nil, err
	}

	acctClient, err := armstorage.NewAccountsClient(mSession.SubscriptionID, aadSession.Cred, aadSession.ClientOptions)
	if err != nil {
		return nil, err
	}
	acctResp, err := acctClient.GetProperties(ctx, resourceGroup, accountName, nil)
	if err != nil {
		return nil, err
	}
	loc := ""
	if acctResp.Account.Location != nil {
		loc = *acctResp.Account.Location
	}
	allowShared := true
	if acctResp.Account.Properties != nil && acctResp.Account.Properties.AllowSharedKeyAccess != nil {
		allowShared = *acctResp.Account.Properties.AllowSharedKeyAccess
	}

	config := GetConfig(d.Connection)
	authMode := "auto"
	if config.DataPlaneAuthMode != nil && *config.DataPlaneAuthMode != "" {
		authMode = strings.ToLower(*config.DataPlaneAuthMode)
	}

	qClient, err := buildQueueServiceClient(ctx, d, aadSession.Cred, authMode, accountName, resourceGroup, mSession.SubscriptionID, mSession.StorageEndpointSuffix, allowShared)
	if err != nil {
		return nil, err
	}

	qClientQueue := qClient.NewQueueClient(name)
	props, err := qClientQueue.GetProperties(ctx, nil)
	if err != nil {
		if queueerror.HasCode(err, queueerror.QueueNotFound) {
			return nil, nil
		}
		return nil, err
	}

	id := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/queueServices/default/queues/%s", mSession.SubscriptionID, resourceGroup, accountName, name)
	meta := map[string]*string{}
	if props.Metadata != nil {
		meta = props.Metadata
	}
	qi := &queueInfo{Name: name, Account: accountName, ResourceGroup: resourceGroup, Location: loc, Metadata: meta, ID: id, Type: "Microsoft.Storage/storageAccounts/queueServices/queues"}
	return qi, nil
}

func buildQueueAkas(_ context.Context, d *transform.TransformData) (interface{}, error) {
	q := d.HydrateItem.(*queueInfo)
	aka := "azure:///subscriptions/" + extractSubscriptionFromID(q.ID) + "/resourceGroups/" + q.ResourceGroup + "/providers/Microsoft.Storage/storageAccounts/" + q.Account + "/queueServices/default/queues/" + q.Name
	akaLower := strings.ToLower(aka)
	return []string{aka, akaLower}, nil
}

// extractSubscriptionFromID quick helper (could reuse existing if available)
func extractSubscriptionFromID(id string) string {
	parts := strings.Split(id, "/")
	for i, p := range parts {
		if p == "subscriptions" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

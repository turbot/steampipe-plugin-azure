select permissions, actions, not_actions, data_actions, not_data_actions
from azure.azure_iam_role_permission
where permissions = '[{"actions":["Microsoft.Authorization/*/read","Microsoft.Compute/*/read","Microsoft.Compute/availabilitySets/*"],"dataActions":["Microsoft.Storage/storageAccounts/blobServices/containers/blobs/delete","Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read","Microsoft.Storage/storageAccounts/blobServices/containers/blobs/write"],"notActions":[],"notDataActions":[]}]'

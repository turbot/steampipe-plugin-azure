select name, storage_account_name, container_name, type, is_snapshot, access_tier, deleted, server_encrypted, content_length, lease_state, region, resource_group, subscription_id
from azure.azure_storage_blob
where title = '{{ resourceName }}'
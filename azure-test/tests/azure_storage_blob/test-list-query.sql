select name, storage_account_name, container_name, type, is_snapshot, access_tier, deleted, server_encrypted, content_length, lease_state, region, resource_group, subscription_id
from azure.azure_storage_blob
where resource_group = '{{ resourceName }}' and storage_account_name = '{{ resourceName }}' and region = '{{ output.location.value }}' and name = '{{ resourceName }}';
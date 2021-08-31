select name, id, type, region, resource_group, subscription_id
from azure.azure_key_vault_deleted_vault
where name = '{{ output.resource_name.value }}' and region = '{{ output.region.value }}';
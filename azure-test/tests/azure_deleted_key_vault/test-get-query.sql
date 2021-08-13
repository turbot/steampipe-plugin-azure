select name, id, type, region, resource_group, subscription_id
from azure.azure_deleted_key_vault
where name = '{{ output.resource_name.value }}' and region = '{{ output.region.value }}';
select name, id, type, kind, sku, region, resource_group, subscription_id, tags
from azure.azure_cognitive_account
where id = '{{ output.resource_id.value }}';

select name, id, type, resource_group, subscription_id
from azure.azure_cdn_frontdoor_profile
where id = '{{ output.resource_id.value }}';

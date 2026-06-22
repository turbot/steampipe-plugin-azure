select name, id, type, profile_name, resource_group, subscription_id
from azure.azure_cdn_frontdoor_custom_domain
where id = '{{ output.resource_id.value }}';

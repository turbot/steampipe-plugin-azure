select name, id, type, profile_name, resource_group, subscription_id
from azure.azure_cdn_frontdoor_custom_domain
where name = '{{ resourceName }}' and profile_name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

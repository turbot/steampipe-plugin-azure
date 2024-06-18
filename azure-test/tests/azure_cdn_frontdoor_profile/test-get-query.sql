select name, id, type, resource_group, subscription_id
from azure.azure_cdn_frontdoor_profile
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

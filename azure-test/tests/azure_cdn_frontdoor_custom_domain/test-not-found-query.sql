select name, id, type
from azure.azure_cdn_frontdoor_custom_domain
where name = 'dummy-test{{ resourceName }}' and profile_name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

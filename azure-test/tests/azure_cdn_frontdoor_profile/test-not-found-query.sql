select name, id, type
from azure.azure_cdn_frontdoor_profile
where name = 'dummy-test{{ resourceName }}' and resource_group = '{{ resourceName }}';

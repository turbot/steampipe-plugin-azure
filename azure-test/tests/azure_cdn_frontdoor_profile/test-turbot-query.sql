select name, akas, title
from azure.azure_cdn_frontdoor_profile
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

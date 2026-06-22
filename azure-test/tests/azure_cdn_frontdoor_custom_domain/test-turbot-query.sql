select name, akas, title
from azure.azure_cdn_frontdoor_custom_domain
where name = '{{ resourceName }}' and profile_name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

select name, akas, title
from azure.azure_frontdoor
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

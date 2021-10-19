select name, id, type
from azure.azure_frontdoor
where name = 'dummy-test{{ resourceName }}' and resource_group = '{{ resourceName }}';

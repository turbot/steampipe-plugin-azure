select name, akas, title
from azure.azure_hpc_cache
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

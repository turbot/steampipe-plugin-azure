select name, akas, tags, title
from azure.azure_storage_sync
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

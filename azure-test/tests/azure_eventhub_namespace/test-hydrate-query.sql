select name, akas, tags, title
from azure.azure_eventhub_namespace
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

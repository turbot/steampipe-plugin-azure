select name, akas, title
from azure.azure_signalr_service
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

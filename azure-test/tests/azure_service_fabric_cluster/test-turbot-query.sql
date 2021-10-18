select name, akas, title
from azure.azure_service_fabric_cluster
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

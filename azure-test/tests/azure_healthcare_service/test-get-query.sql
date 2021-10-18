select name, id, kind, region, type
from azure.azure_healthcare_service
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
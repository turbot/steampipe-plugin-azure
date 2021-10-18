select name, akas, title, tags
from azure.azure_healthcare_service
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

select name, id
from azure.azure_healthcare_service
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';

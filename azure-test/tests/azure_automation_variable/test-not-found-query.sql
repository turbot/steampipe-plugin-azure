select name, title, akas
from azure.azure_automation_variable
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';
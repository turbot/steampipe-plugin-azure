select name, tags, title, akas
from azure.azure_automation_account
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';
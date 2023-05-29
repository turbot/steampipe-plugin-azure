select name, title, akas
from azure.azure_automation_account
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
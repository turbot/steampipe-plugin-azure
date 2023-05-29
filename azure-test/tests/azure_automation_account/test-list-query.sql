select name, id, type
from azure.azure_automation_account
where name = '{{ resourceName }}';
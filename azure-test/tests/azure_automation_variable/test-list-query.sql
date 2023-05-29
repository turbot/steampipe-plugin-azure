select name, id, type
from azure.azure_automation_variable
where name = '{{ resourceName }}';
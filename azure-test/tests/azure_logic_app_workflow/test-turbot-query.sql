select name, title, akas
from azure.azure_logic_app_workflow
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
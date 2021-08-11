select name, tags, title, akas
from azure.azure_logic_app_workflow
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';
select name, tags, title, akas
from azure.azure_machine_learning_workspace
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';
select name, tags, title, akas
from azure.azure_machine_learning_workspace
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
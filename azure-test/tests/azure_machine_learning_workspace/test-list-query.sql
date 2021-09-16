select name, id, type
from azure.azure_machine_learning_workspace
where name = '{{ resourceName }}';
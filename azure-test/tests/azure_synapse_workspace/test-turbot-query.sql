select name, akas, title, tags
from azure.azure_synapse_workspace
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

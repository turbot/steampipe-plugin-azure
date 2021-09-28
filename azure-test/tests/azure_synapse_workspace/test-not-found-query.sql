select name, id, type, region
from azure.azure_synapse_workspace
where name = 'dummy-test{{ resourceName }}' and resource_group = '{{ resourceName }}';

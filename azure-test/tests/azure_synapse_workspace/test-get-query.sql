select name, id, type, region, resource_group, subscription_id, tags, sql_administrator_login
from azure.azure_synapse_workspace
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

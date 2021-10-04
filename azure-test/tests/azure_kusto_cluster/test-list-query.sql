select id, name
from azure.azure_kusto_cluster
where name = '{{ resourceName }}';

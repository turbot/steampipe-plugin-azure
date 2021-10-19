select name, akas, title
from azure.azure_kusto_cluster
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

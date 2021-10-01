select name, id, sku_name, uri
from azure.azure_kusto_cluster
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

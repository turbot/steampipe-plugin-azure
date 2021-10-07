select name, id, type, region, resource_group, subscription_id
from azure.azure_hdinsight_cluster
where id = '{{ output.resource_id.value }}';

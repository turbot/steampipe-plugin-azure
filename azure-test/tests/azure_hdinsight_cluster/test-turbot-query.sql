select name, akas, title
from azure.azure_hdinsight_cluster
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

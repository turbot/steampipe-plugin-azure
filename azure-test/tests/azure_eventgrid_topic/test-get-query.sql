select name, id, type
from azure.azure_eventgrid_topic
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';

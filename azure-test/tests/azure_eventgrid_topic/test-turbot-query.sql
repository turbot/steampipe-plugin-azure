select name, akas, title
from azure.azure_eventgrid_topic
where name = '{{ resourceName }}';

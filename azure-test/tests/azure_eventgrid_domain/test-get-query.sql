select name, id, type
from azure.azure_eventgrid_domain
where name = '{{ resourceName }}';

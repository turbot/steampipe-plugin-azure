select name, akas, title
from azure.azure_eventgrid_domain
where name = '{{ resourceName }}';

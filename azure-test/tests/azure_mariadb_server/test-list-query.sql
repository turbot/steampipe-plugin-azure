select
  name,
  id
from
  azure.azure_mariadb_server
where
  id = '{{ output.resource_id.value }}';
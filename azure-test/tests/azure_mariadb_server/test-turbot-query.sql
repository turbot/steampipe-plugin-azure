select
  name,
  id,
  region,
  resource_group,
  subscription_id,
  title
from
  azure.azure_mariadb_server
where
  name = '{{ resourceName }}'
  and resource_group = '{{ resourceName }}';
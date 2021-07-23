select
  name,
  id,
  region,
  resource_group,
  subscription_id,
  title
from
  azure.azure_container_registry
where
  name = '{{ resourceName }}'
  and resource_group = '{{ resourceName }}';
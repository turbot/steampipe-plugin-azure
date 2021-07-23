select
  name,
  id
from
  azure.azure_container_registry
where
  name = 'dummy-{{ resourceName }}'
  and resource_group = '{{ resourceName }}';
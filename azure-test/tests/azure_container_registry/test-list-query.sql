select
  name,
  id
from
  azure.azure_container_registry
where
  id = '{{ output.resource_id.value }}';
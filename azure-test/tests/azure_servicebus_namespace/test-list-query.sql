select
  name,
  id
from
  azure.azure_servicebus_namespace
where
  id = '{{ output.resource_id.value }}';
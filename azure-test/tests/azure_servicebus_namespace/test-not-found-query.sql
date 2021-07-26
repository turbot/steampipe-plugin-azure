select
  name,
  id
from
  azure.azure_servicebus_namespace
where
  name = 'dummy-{{ resourceName }}'
  and resource_group = '{{ resourceName }}';
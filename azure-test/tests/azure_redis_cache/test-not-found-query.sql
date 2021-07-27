select
  name,
  id
from
  azure.azure_redis_cache
where
  name = 'dummy-{{ resourceName }}'
  and resource_group = '{{ resourceName }}';
select
  name,
  type
from
  azure.azure_redis_cache
where
  name = '{{ resourceName }}';
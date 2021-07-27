select
  name,
  tags,
  title,
  akas
from
  azure.azure_redis_cache
where
  name = '{{ resourceName }}'
  and resource_group = '{{ resourceName }}';
select
  name,
  type,
  redis_version,
  enable_non_ssl_port,
  host_name,
  minimum_tls_version,
  public_network_access,
  sku_capacity,
  sku_family,
  sku_name,
  region,
  resource_group,
  subscription_id
from
  azure.azure_redis_cache
where
  name = '{{ resourceName }}'
  and resource_group = '{{ resourceName }}';
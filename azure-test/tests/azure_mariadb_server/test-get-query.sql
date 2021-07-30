select
  name,
  id,
  type,
  version,
  geo_redundant_backup_enabled,
  user_visible_state,
  administrator_login,
  auto_grow_enabled,
  backup_retention_days,
  fully_qualified_domain_name,
  public_network_access,
  sku_capacity,
  sku_family,
  sku_name,
  sku_tier,
  ssl_enforcement,
  storage_mb
from
  azure.azure_mariadb_server
where
  name = '{{ resourceName }}'
  and resource_group = '{{ resourceName }}';
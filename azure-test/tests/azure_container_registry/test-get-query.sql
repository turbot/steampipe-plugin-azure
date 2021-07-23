select
  name,
  id,
  type,
  provisioning_state,
  admin_user_enabled,
  data_endpoint_enabled,
  login_server,
  network_rule_bypass_options,
  public_network_access,
  sku_name,
  sku_tier,
  zone_redundancy,
  encryption,
  policies
from
  azure.azure_container_registry
where
  name = '{{ resourceName }}'
  and resource_group = '{{ resourceName }}';
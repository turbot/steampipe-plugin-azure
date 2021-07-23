select
  name,
  id,
  type,
  provisioning_state,
  zone_redundant,
  servicebus_endpoint,
  sku_name,
  sku_tier,
  encryption,
  network_rule_set
from
  azure.azure_servicebus_namespace
where
  name = '{{ resourceName }}'
  and resource_group = '{{ resourceName }}';
select
  name,
  id,
  title,
  region,
  resource_group,
  subscription_id,
  tags,
  akas
from
  azure.azure_servicebus_namespace
where
  name = '{{ resourceName }}'
  and resource_group = '{{ resourceName }}';
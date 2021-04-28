select
  name,
  id,
  vault_name,
  enabled,
  content_type,
  recoverable_days,
  recovery_level,
  value,
  region,
  resource_group,
  subscription_id
from
  azure.azure_key_vault_secret
where
  name = '{{ resourceName }}'
  and vault_name = '{{ resourceName }}';
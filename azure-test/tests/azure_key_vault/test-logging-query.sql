select
  name,
  setting -> 'properties' ->> 'storageAccountId' storage_account_id,
  log ->> 'category' category,
  (log -> 'retentionPolicy' ->> 'days')::integer log_retention_days
from
  azure_key_vault,
  jsonb_array_elements(diagnostic_settings) setting,
  jsonb_array_elements(setting -> 'properties' -> 'logs') log
where
  diagnostic_settings is not null
  and setting -> 'properties' ->> 'storageAccountId' <> ''
  and (log ->> 'enabled')::boolean
  and log ->> 'category' = 'AuditEvent'
  and (log -> 'retentionPolicy' ->> 'days')::integer > 0;
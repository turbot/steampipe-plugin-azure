select
  name,
  id,
  storage_account_id
from
  azure.azure_monitor_log_profile
where
  name = '{{ resourceName }}'
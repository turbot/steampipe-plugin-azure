select
  name,
  subscription_id,
  title
from
  azure.azure_monitor_log_profile
where
  name = '{{ resourceName }}';
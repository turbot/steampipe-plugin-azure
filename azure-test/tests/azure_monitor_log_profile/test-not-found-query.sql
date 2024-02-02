select
  name,
  id
from
  azure.azure_monitor_log_profile
where
  name = 'dummy-{{ resourceName }}';
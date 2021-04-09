select 
  name, 
  id,
  account_name,
  deleted,
  public_access,
  deny_encryption_scope_override,
  has_immutability_policy,
  has_legal_hold,
  lease_status,
  lease_state,
  lease_duration,
  remaining_retention_days,
  version,
  type,
  resource_group,
  subscription_id
from 
  azure.azure_storage_container
where 
  name = '{{ resourceName }}' 
  and resource_group = '{{ resourceName }}' 
  and  account_name = '{{ resourceName }}';
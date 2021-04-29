select 
  id, 
  name
from 
  azure.azure_key_vault_key
where 
  name = '{{ resourceName }}';

select 
  id, 
  name
from 
  azure.azure_key_vault_secret
where 
  name = '{{ resourceName }}';

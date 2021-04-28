select 
  name, 
  akas, 
  title, 
  tags
from 
  azure.azure_key_vault_secret
where 
  name = '{{ resourceName }}' 
  and vault_name = '{{ resourceName }}';

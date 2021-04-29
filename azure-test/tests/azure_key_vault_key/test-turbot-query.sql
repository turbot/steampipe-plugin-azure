select 
  name, 
  akas, 
  title, 
  tags
from 
  azure.azure_key_vault_key
where 
  name = '{{ resourceName }}' 
  and resource_group = '{{ resourceName }}' 
  and vault_name = '{{ resourceName }}';

# Table: azure_compute_ssh_key

Azure SSH public key used by VMs. 

## Examples

### Retrieve SSH public key by name

```sql
select
  name,
  publicKey
from
    table_azure_compute_ssh_key
where
  name = 'key-name.';
```

### List compute virtual machines using SSH public key

```sql
select
  m.name as machine_name,
  k.name as ssh_key_name
from
  azure_compute_virtual_machine as m,
  jsonb_array_elements(linux_configuration_ssh_public_keys) as s
  left join azure_compute_ssh_key as k on k.public_key = s ->> 'keyData';
```

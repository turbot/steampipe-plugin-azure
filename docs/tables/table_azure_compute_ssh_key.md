# Table: table_azure_compute_ssh_key

Azure SSH public key used by VM. 

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


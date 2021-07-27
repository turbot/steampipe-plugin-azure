# Table: azure_security_center_jit_network_access_policy

Azure Security Center's just-in-time (JIT) network access policy helps to lock down inbound traffic access to your virtual machines. This reduces exposure to attacks while providing easy access when you need to connect to a VM.

## Examples

### List virtual machines with JIT access enabled

```sql
select
  vm.name,
  vm.id,
  jsonb_pretty(vms -> 'ports') as ports
from
  azure_security_center_jit_network_access_policy,
  jsonb_array_elements(virtual_machines) as vms,
  azure_compute_virtual_machine as vm
where
  lower(vms ->> 'id') = lower(vm.id);
```

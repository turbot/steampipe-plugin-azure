# Table: azure_compute_virtual_machine_scale_set_network_interface

A network interface enables an Azure VM to communicate with internet, Azure, and on-premises resources.

## Examples

### Basic info

```sql
select
  name,
  id,
  scale_set_name,
  provisioning_state,
  region,
  resource_group
from
  azure_compute_virtual_machine_scale_set_network_interface;
```

### List network interfaces with IP forwarding rule enabled

```sql
select
  name,
  id,
  enable_ip_forwarding,
  sku_name
from
  azure_compute_virtual_machine_scale_set_network_interface
where
  enable_ip_forwarding;
```

### List network interfaces with accelerated networking enabled

```sql
select
  name,
  id,
  enable_accelerated_networking,
  sku_name
from
  azure_compute_virtual_machine_scale_set_network_interface
where
  enable_accelerated_networking;
```

### Get scale set virtual machine details for scale set network interface

```sql
select
  i.name as name,
  i.id as id,
  v.instance_id as instance_id,
  v.scale_set_name as scale_set_name,
  v.sku_name as vm_sku_name
from
  azure_compute_virtual_machine_scale_set_network_interface as i,
  azure_compute_virtual_machine_scale_set_vm as v
where
  i.virtual_machine ->> 'id' = v.id;
```

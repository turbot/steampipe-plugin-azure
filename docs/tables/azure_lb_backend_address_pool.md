# Table: azure_lb_backend_address_pool

An Azure load balancer's backend address pool consists of IP addresses associated with the virtual machine NICs. This pool is used to distribute traffic to the virtual machines behind the load balancer.

## Examples

### Basic info

```sql
select
  id,
  name,
  load_balancer_name,
  provisioning_state,
  type
from
  azure_lb_backend_address_pool;
```

### List failed load balancer backend address pools

```sql
select
  id,
  name,
  type,
  provisioning_state
from
  azure_lb_backend_address_pool
where
  provisioning_state = 'Failed';
```

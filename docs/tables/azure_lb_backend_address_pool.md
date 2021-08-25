# Table: azure_lb_backend_address_pool

The backend pool is a critical component of the load balancer. The backend pool defines the group of resources that will serve traffic for a given load-balancing rule.

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

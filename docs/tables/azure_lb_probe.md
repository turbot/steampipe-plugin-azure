# Table: azure_lb_probe

A health probe is used to determine the health status of the instances in the backend pool. It will determine if an instance is healthy and can receive traffic.

## Examples

### Basic info

```sql
select
  name,
  type,
  provisioning_state,
  load_balancer_name,
  port
from
  azure_lb_probe;
```

### List succeeded load balancer probe

```sql
select
  name,
  provisioning_state
from
  azure_lb_probe
where
  provisioning_state = 'Succeeded'
```

### List load balancer probe order by interval

```sql
select
  name,
  interval_in_seconds
from
  azure_lb_probe
order by 
  interval_in_seconds
```
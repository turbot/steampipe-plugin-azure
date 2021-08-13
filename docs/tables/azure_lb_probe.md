# Table: azure_lb_probe

When using load-balancing rules with Azure Load Balancer, you need to specify health probes to allow Load Balancer to detect the backend endpoint status. The configuration of the health probe and probe responses determine which backend pool instances will receive new flows. You can use health probes to detect the failure of an application on a backend endpoint.

## Examples

### Basic info

```sql
select
  id,
  name,
  type,
  provisioning_state,
  load_balancer_name,
  port
from
  azure_lb_probe;
```

### List failed load balancer probes

```sql
select
  id,
  name,
  type,
  provisioning_state
from
  azure_lb_probe
where
  provisioning_state = 'Failed';
```

### List load balancer probes order by interval

```sql
select
  id,
  name,
  type,
  interval_in_seconds
from
  azure_lb_probe
order by 
  interval_in_seconds;
```
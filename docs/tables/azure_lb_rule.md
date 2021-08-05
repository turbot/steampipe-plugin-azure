# Table: azure_lb_rule

A load balancer rule is used to define how traffic is distributed to the VMs. You define the front-end IP configuration for the incoming traffic and the back-end IP pool to receive the traffic, along with the required source and destination port.

## Examples

### Basic info

```sql
select
  name,
  type,
  provisioning_state,
  etag
from
  azure_lb_rule;
```

### List load balancer rule with succeeded provisioning state

```sql
select
  name,
  provisioning_state
from
  azure_lb_rule
where
  provisioning_state = 'Succeeded'
```

### Order load balancer rule by idle timeout

```sql
select
  name,
  idle_timeout_in_minutes
from
  azure_lb_rule
order by 
  idle_timeout_in_minutes
```
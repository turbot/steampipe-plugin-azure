# Table: azure_lb_nat_rule

An inbound NAT rule forwards incoming traffic sent to frontend IP address and port combination. The traffic is sent to a specific virtual machine or instance in the backend pool. Port forwarding is done by the same hash-based distribution as load balancing.

## Examples

### Basic info

```sql
select
  id,
  name,
  type,
  provisioning_state,
  etag
from
  azure_lb_nat_rule;
```

### List failed load balancer nat rules

```sql
select
  id,
  name,
  type,
  provisioning_state
from
  azure_lb_nat_rule
where
  provisioning_state = 'Failed';
```

### List load balancer nat rules order by idle timeout

```sql
select
  id,
  name,
  type,
  idle_timeout_in_minutes
from
  azure_lb_nat_rule
order by 
  idle_timeout_in_minutes;
```
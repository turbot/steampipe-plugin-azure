# Table: azure_lb_outbound_rule

Outbound rules allow you to explicitly define SNAT(source network address translation) for a public standard load balancer. This configuration allows you to use the public IP(s) of your load balancer to provide outbound internet connectivity for your backend instances.

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
  azure_lb_outbound_rule;
```

### List failed load balancer outbound rules

```sql
select
  id,
  name,
  type,
  provisioning_state
from
  azure_lb_outbound_rule
where
  provisioning_state = 'Failed';
```

### List load balancer outbound rules order by idle timeout

```sql
select
  id,
  name,
  type,
  idle_timeout_in_minutes
from
  azure_lb_outbound_rule
order by 
  idle_timeout_in_minutes;
```
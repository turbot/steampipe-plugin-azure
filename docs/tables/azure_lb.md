# Table: azure_lb

Azure Load Balancer operates at layer 4 of the Open Systems Interconnection (OSI) model. It's the single point of contact for clients. Load balancer distributes inbound flows that arrive at the load balancer's front end to backend pool instances. These flows are according to configured load-balancing rules and health probes. The backend pool instances can be Azure Virtual Machines or instances in a virtual machine scale set.

## Examples

### Basic info

```sql
select
  id,
  name,
  type,
  provisioning_state,
  etag,
  region
from
  azure_lb;
```

### List failed load balancers

```sql
select
  id,
  name,
  type,
  provisioning_state
from
  azure_lb
where
  provisioning_state = 'Failed';
```
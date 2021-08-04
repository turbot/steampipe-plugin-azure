# Table: azure_lb

An Azure load balancer is a Layer-4 (TCP, UDP) load balancer that provides high availability by distributing incoming traffic among healthy VMs.

## Examples

### Basic info

```sql
select
  name,
  type,
  provisioning_state,
  etag,
  region
from
  azure_lb;
```

### List load balancer with failed provisioning state

```sql
select
  name,
  provisioning_state
from
  azure_lb
where
  provisioning_state = 'Failed'
```
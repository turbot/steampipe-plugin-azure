# Table: azure_security_center_subscription_pricing

Azure security center pricing configurations for the subscription.

## Examples

### Basic info

```sql
select
  id,
  name,
  pricing_tier
from
  azure_security_center_subscription_pricing;
```

### Get pricing tier for VirtualMachines

```sql
select
  id,
  name,
  pricing_tier
from
  azure_security_center_subscription_pricing
where
  name = 'VirtualMachines';
```

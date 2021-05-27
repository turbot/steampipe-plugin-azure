# Table: azure_tenant

A dedicated and trusted instance of Azure AD that's automatically created when your organization signs up for a Microsoft cloud service subscription, such as Microsoft Azure, Microsoft Intune, or Microsoft 365. An Azure tenant represents a single organization.

## Examples

### Basic info

```sql
select
  name,
  id,
  tenant_id,
  tenant_category,
  country,
  country_code,
  display_name,
  domains
from
  azure_tenant;
```

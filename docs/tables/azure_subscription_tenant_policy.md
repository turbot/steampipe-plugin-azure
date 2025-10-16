---
title: "Steampipe Table: azure_subscription_tenant_policy - Query Azure Subscription Tenant Policies using SQL"
description: "Allows users to query Azure Subscription Tenant Policies, providing insights into tenant-level policies that control subscription movement between tenants."
folder: "Subscription"
---

# Table: azure_subscription_tenant_policy - Query Azure Subscription Tenant Policies using SQL

Azure Subscription Tenant Policies control whether subscriptions can leave or enter an Azure tenant. These policies help organizations maintain control over their subscription governance by restricting the movement of subscriptions across tenant boundaries. Tenant administrators can configure these policies to prevent unauthorized subscription transfers and maintain organizational compliance.

## Table Usage Guide

The `azure_subscription_tenant_policy` table provides insights into tenant-level subscription policies within Microsoft Azure. As a tenant administrator or cloud governance specialist, explore policy details through this table, including settings that block subscriptions from leaving or entering the tenant, and lists of exempted principals. Utilize it to monitor and enforce subscription governance policies, ensuring compliance with organizational security and management requirements.

## Examples

### Basic info
Explore the tenant policy settings to understand which restrictions are in place for subscription movement within your Azure tenant.

```sql+postgres
select
  name,
  id,
  policy_id,
  block_subscriptions_leaving_tenant,
  block_subscriptions_into_tenant
from
  azure_subscription_tenant_policy;
```

```sql+sqlite
select
  name,
  id,
  policy_id,
  block_subscriptions_leaving_tenant,
  block_subscriptions_into_tenant
from
  azure_subscription_tenant_policy;
```

### Check if subscriptions are blocked from leaving the tenant
Determine whether your tenant has policies in place to prevent subscriptions from being transferred to other tenants.

```sql+postgres
select
  name,
  policy_id,
  block_subscriptions_leaving_tenant,
  case
    when block_subscriptions_leaving_tenant then 'Restricted'
    else 'Allowed'
  end as leaving_status
from
  azure_subscription_tenant_policy;
```

```sql+sqlite
select
  name,
  policy_id,
  block_subscriptions_leaving_tenant,
  case
    when block_subscriptions_leaving_tenant then 'Restricted'
    else 'Allowed'
  end as leaving_status
from
  azure_subscription_tenant_policy;
```

### Check if subscriptions are blocked from entering the tenant
Identify whether your tenant has restrictions on accepting subscriptions from other tenants.

```sql+postgres
select
  name,
  policy_id,
  block_subscriptions_into_tenant,
  case
    when block_subscriptions_into_tenant then 'Restricted'
    else 'Allowed'
  end as entering_status
from
  azure_subscription_tenant_policy;
```

```sql+sqlite
select
  name,
  policy_id,
  block_subscriptions_into_tenant,
  case
    when block_subscriptions_into_tenant then 'Restricted'
    else 'Allowed'
  end as entering_status
from
  azure_subscription_tenant_policy;
```

### List exempted principals
Discover which user principals are exempted from the subscription tenant policies in your organization.

```sql+postgres
select
  name,
  policy_id,
  jsonb_array_elements_text(exempted_principals) as exempted_principal
from
  azure_subscription_tenant_policy;
```

```sql+sqlite
select
  name,
  policy_id,
  json_each.value as exempted_principal
from
  azure_subscription_tenant_policy,
  json_each(exempted_principals);
```

### Get tenant policies with both entering and leaving restrictions
Find tenant policies that have comprehensive restrictions preventing both incoming and outgoing subscription transfers.

```sql+postgres
select
  name,
  policy_id,
  block_subscriptions_leaving_tenant,
  block_subscriptions_into_tenant
from
  azure_subscription_tenant_policy
where
  block_subscriptions_leaving_tenant
  and block_subscriptions_into_tenant;
```

```sql+sqlite
select
  name,
  policy_id,
  block_subscriptions_leaving_tenant,
  block_subscriptions_into_tenant
from
  azure_subscription_tenant_policy
where
  block_subscriptions_leaving_tenant = 1
  and block_subscriptions_into_tenant = 1;
```


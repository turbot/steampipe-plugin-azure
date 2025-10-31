---
title: "Steampipe Table: azure_subscription - Query Azure Subscriptions using SQL"
description: "Allows users to query Azure Subscriptions, providing insights into subscription details, including subscription IDs, names, states, and tenants."
folder: "Subscription"
---

# Table: azure_subscription - Query Azure Subscriptions using SQL

Azure Subscriptions act as a logical container for resources deployed on Microsoft Azure. They provide a mechanism to organize access to Azure resources, manage costs, and track billing. Each Azure Subscription can have a different billing and payment setup, allowing flexibility in how users and organizations pay for the usage of Azure Services.

## Table Usage Guide

The `azure_subscription` table provides insights into Azure Subscriptions within Microsoft Azure. As a cloud architect or administrator, explore subscription-specific details through this table, including subscription IDs, names, states, and tenants. Utilize it to manage and organize access to Azure resources, track billing, and understand the cost management setup across different subscriptions.

## Examples

### Basic info
Explore the status and policies of your Azure subscriptions to understand their current state and source of authorization. This can help in managing and optimizing your cloud resources effectively.

```sql+postgres
select
  id,
  subscription_id,
  display_name,
  tenant_id,
  state,
  authorization_source,
  subscription_policies
from
  azure_subscription;
```

```sql+sqlite
select
  id,
  subscription_id,
  display_name,
  tenant_id,
  state,
  authorization_source,
  subscription_policies
from
  azure_subscription;
```

### Get tenant policy settings for subscriptions
Retrieve the tenant policy configuration to understand restrictions on subscription transfers. This helps identify whether subscriptions are blocked from leaving or entering the tenant, and which principals are exempted from these policies.

```sql+postgres
select
  subscription_id,
  display_name,
  default_tenant_policy ->> 'id' as policy_id,
  default_tenant_policy -> 'properties' ->> 'blockSubscriptionsLeavingTenant' as block_leaving_tenant,
  default_tenant_policy -> 'properties' ->> 'blockSubscriptionsIntoTenant' as block_into_tenant,
  default_tenant_policy -> 'properties' -> 'exemptedPrincipals' as exempted_principals
from
  azure_subscription;
```

```sql+sqlite
select
  subscription_id,
  display_name,
  json_extract(default_tenant_policy, '$.id') as policy_id,
  json_extract(default_tenant_policy, '$.properties.blockSubscriptionsLeavingTenant') as block_leaving_tenant,
  json_extract(default_tenant_policy, '$.properties.blockSubscriptionsIntoTenant') as block_into_tenant,
  json_extract(default_tenant_policy, '$.properties.exemptedPrincipals') as exempted_principals
from
  azure_subscription;
```
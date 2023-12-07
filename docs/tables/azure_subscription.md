---
title: "Steampipe Table: azure_subscription - Query Azure Subscriptions using SQL"
description: "Allows users to query Azure Subscriptions, providing insights into subscription details, including subscription IDs, names, states, and tenants."
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
---
title: "Steampipe Table: azure_subscription - Query Azure Subscriptions using SQL"
description: "Allows users to query Azure Subscriptions."
---

# Table: azure_subscription - Query Azure Subscriptions using SQL

Azure Subscriptions represent a logical container for resources that are deployed within an Azure account. They provide a way to manage costs and resources where users can apply different policies and manage access control. Each Azure subscription can have a separate billing and payment setup, so you can have different subscriptions for different departments or projects.

## Table Usage Guide

The 'azure_subscription' table provides insights into subscriptions within Azure. As a DevOps engineer, explore subscription-specific details through this table, including subscription ID, name, and state, among others. Utilize it to uncover information about subscriptions, such as their current state, the tenant they belong to, and whether they are spending over their budget. The schema presents a range of attributes of the Azure subscription for your analysis, like the subscription ID, tenant ID, state, and location placement ID.

## Examples

### Basic info
Explore which Azure subscriptions are active and the policies associated with them. This can be helpful in managing resources and understanding the scope of your Azure environment.

```sql
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
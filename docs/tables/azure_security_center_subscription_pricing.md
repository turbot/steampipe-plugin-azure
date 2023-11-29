---
title: "Steampipe Table: azure_security_center_subscription_pricing - Query Azure Security Center Subscription Pricings using SQL"
description: "Allows users to query Azure Security Center Subscription Pricings."
---

# Table: azure_security_center_subscription_pricing - Query Azure Security Center Subscription Pricings using SQL

Azure Security Center is a unified security management system that strengthens the security posture of your data centers and provides advanced threat protection across your hybrid workloads in the cloud, whether they're in Azure or not. It provides you with a set of policies and recommendations, tailored to your specific deployments. With Azure Security Center, you can understand the security state of your resources in Azure, on-premises, and in other cloud providers.

## Table Usage Guide

The 'azure_security_center_subscription_pricing' table provides insights into subscription pricings within Azure Security Center. As a security analyst, explore pricing-specific details through this table, including pricing tier, free trial status, and associated metadata. Utilize it to uncover information about subscription pricings, such as the pricing tier for each resource type and whether the free trial is still active. The schema presents a range of attributes of the subscription pricing for your analysis, like the pricing name, pricing tier, and free trial status.

## Examples

### Basic info
Analyze the settings to understand the different pricing tiers of your Azure Security Center subscriptions. This can help you assess your current cost structure and potentially identify areas for optimization.

```sql
select
  id,
  name,
  pricing_tier
from
  azure_security_center_subscription_pricing;
```

### List pricing information for virtual machines
Explore the cost implications of your virtual machines by determining their associated pricing tiers. This is useful for budget management and cost optimization within your Azure environment.

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
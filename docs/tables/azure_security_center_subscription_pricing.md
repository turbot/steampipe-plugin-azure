---
title: "Steampipe Table: azure_security_center_subscription_pricing - Query Azure Security Center Subscription Pricing using SQL"
description: "Allows users to query Azure Security Center Subscription Pricing, specifically the pricing tier and the resource details associated with the subscription."
---

# Table: azure_security_center_subscription_pricing - Query Azure Security Center Subscription Pricing using SQL

Azure Security Center is a unified infrastructure security management system that strengthens the security posture of your data centers and provides advanced threat protection across your hybrid workloads in the cloud. It provides you with a comprehensive view of your security state and actionable recommendations to mitigate risks. The subscription pricing model allows you to choose the level of protection that best meets your needs.

## Table Usage Guide

The `azure_security_center_subscription_pricing` table provides insights into the pricing tier and resource details associated with each Azure Security Center subscription. As a security analyst, use this table to understand the cost implications of your security strategies, and to ensure you are utilizing the most appropriate level of protection for your needs. This table can also assist in budget planning and cost management for your Azure resources.

## Examples

### Basic info
Explore which Azure Security Center subscriptions are at different pricing tiers to manage costs effectively and ensure optimal resource utilization.

```sql+postgres
select
  id,
  name,
  pricing_tier
from
  azure_security_center_subscription_pricing;
```

```sql+sqlite
select
  id,
  name,
  pricing_tier
from
  azure_security_center_subscription_pricing;
```

### List pricing information for virtual machines
Explore the cost implications of your virtual machines by examining their pricing tiers. This allows for efficient budget management and cost optimization.

```sql+postgres
select
  id,
  name,
  pricing_tier
from
  azure_security_center_subscription_pricing
where
  name = 'VirtualMachines';
```

```sql+sqlite
select
  id,
  name,
  pricing_tier
from
  azure_security_center_subscription_pricing
where
  name = 'VirtualMachines';
```
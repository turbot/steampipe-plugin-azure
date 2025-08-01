---
title: "Steampipe Table: azure_cost_usage - Query Azure Cost and Usage Data using SQL"
description: "Allows users to query Azure Cost and Usage Data with flexible dimensions, providing detailed cost breakdown by any combination of Azure dimensions."
folder: "Cost Management"
---

# Table: azure_cost_usage - Query Azure Cost and Usage Data using SQL

Azure Cost Management provides cost analytics to help you understand and manage your Azure spending. The cost usage table allows you to query cost data with flexible dimensions, providing insights into Azure costs broken down by any combination of supported dimensions such as service name, resource group, location, resource type, subscription, and more. This enables custom cost analysis tailored to your specific organizational needs.

## Table Usage Guide

The `azure_cost_usage` table provides insights into cost and usage data within Microsoft Azure with flexible dimension support. As a Cloud Architect, FinOps engineer, or DevOps professional, explore cost details through this table using any combination of Azure dimensions. Utilize it to create custom cost breakdowns, analyze spending patterns across multiple dimensions, track costs by location and service, and perform advanced cost analytics that match your organizational structure.

**Note:** This table requires three key qualifiers: `granularity` (DAILY or MONTHLY), `dimension_type_1`, and `dimension_type_2`. Supported dimension types include: ResourceGroup, ResourceGroupName, ResourceLocation, ConsumedService, ResourceType, ServiceName, SubscriptionName, MeterCategory, and many others.

## Examples

### Basic cost breakdown by service and resource group

Explore daily costs broken down by service name and resource group to understand which services are costing the most in each resource group.

```sql+postgres
select
  usage_date,
  dimension_1 as service_name,
  dimension_2 as resource_group,
  unblended_cost_amount,
  unblended_cost_unit
from
  azure_cost_usage
where
  granularity = 'DAILY'
  and dimension_type_1 = 'ServiceName'
  and dimension_type_2 = 'ResourceGroupName'
order by
  usage_date desc,
  unblended_cost_amount desc;
```

```sql+sqlite
select
  usage_date,
  dimension_1 as service_name,
  dimension_2 as resource_group,
  unblended_cost_amount,
  unblended_cost_unit
from
  azure_cost_usage
where
  granularity = 'DAILY'
  and dimension_type_1 = 'ServiceName'
  and dimension_type_2 = 'ResourceGroupName'
order by
  usage_date desc,
  unblended_cost_amount desc;
```

### Monthly costs by location and service

Analyze monthly costs broken down by resource location and service to understand geographical spending patterns.

```sql+postgres
select
  usage_date,
  dimension_1 as location,
  dimension_2 as service_name,
  unblended_cost_amount,
  unblended_cost_unit
from
  azure_cost_usage
where
  granularity = 'MONTHLY'
  and dimension_type_1 = 'ResourceLocation'
  and dimension_type_2 = 'ServiceName'
order by
  usage_date desc,
  unblended_cost_amount desc;
```

```sql+sqlite
select
  usage_date,
  dimension_1 as location,
  dimension_2 as service_name,
  unblended_cost_amount,
  unblended_cost_unit
from
  azure_cost_usage
where
  granularity = 'MONTHLY'
  and dimension_type_1 = 'ResourceLocation'
  and dimension_type_2 = 'ServiceName'
order by
  usage_date desc,
  unblended_cost_amount desc;
```

### Resource type and consumed service analysis

Understand costs by resource type and the services that consume them for detailed infrastructure cost analysis.

```sql+postgres
select
  usage_date,
  dimension_1 as resource_type,
  dimension_2 as consumed_service,
  unblended_cost_amount,
  unblended_cost_unit
from
  azure_cost_usage
where
  granularity = 'DAILY'
  and dimension_type_1 = 'ResourceType'
  and dimension_type_2 = 'ConsumedService'
  and usage_date >= current_date - interval '7 days'
order by
  unblended_cost_amount desc;
```

```sql+sqlite
select
  usage_date,
  dimension_1 as resource_type,
  dimension_2 as consumed_service,
  unblended_cost_amount,
  unblended_cost_unit
from
  azure_cost_usage
where
  granularity = 'DAILY'
  and dimension_type_1 = 'ResourceType'
  and dimension_type_2 = 'ConsumedService'
  and usage_date >= date('now', '-7 days')
order by
  unblended_cost_amount desc;
```

### Subscription and service cost breakdown

Track costs across different subscriptions and services to understand multi-subscription spending patterns.

```sql+postgres
select
  usage_date,
  dimension_1 as subscription_name,
  dimension_2 as service_name,
  sum(unblended_cost_amount) as total_cost,
  unblended_cost_unit
from
  azure_cost_usage
where
  granularity = 'DAILY'
  and dimension_type_1 = 'SubscriptionName'
  and dimension_type_2 = 'ServiceName'
  and usage_date >= current_date - interval '30 days'
group by
  usage_date,
  dimension_1,
  dimension_2,
  unblended_cost_unit
order by
  usage_date desc,
  total_cost desc;
```

```sql+sqlite
select
  usage_date,
  dimension_1 as subscription_name,
  dimension_2 as service_name,
  sum(unblended_cost_amount) as total_cost,
  unblended_cost_unit
from
  azure_cost_usage
where
  granularity = 'DAILY'
  and dimension_type_1 = 'SubscriptionName'
  and dimension_type_2 = 'ServiceName'
  and usage_date >= date('now', '-30 days')
group by
  usage_date,
  dimension_1,
  dimension_2,
  unblended_cost_unit
order by
  usage_date desc,
  total_cost desc;
```

### Top cost combinations for specific date range

Find the highest cost combinations for any dimension pair within a specific date range.

```sql+postgres
select
  dimension_1,
  dimension_2,
  sum(unblended_cost_amount) as total_cost,
  count(*) as usage_days,
  unblended_cost_unit
from
  azure_cost_usage
where
  granularity = 'DAILY'
  and dimension_type_1 = 'ServiceName'
  and dimension_type_2 = 'ResourceGroupName'
  and usage_date between '2025-01-01' and '2025-01-31'
group by
  dimension_1,
  dimension_2,
  unblended_cost_unit
order by
  total_cost desc
limit 10;
```

```sql+sqlite
select
  dimension_1,
  dimension_2,
  sum(unblended_cost_amount) as total_cost,
  count(*) as usage_days,
  unblended_cost_unit
from
  azure_cost_usage
where
  granularity = 'DAILY'
  and dimension_type_1 = 'ServiceName'
  and dimension_type_2 = 'ResourceGroupName'
  and usage_date between '2025-01-01' and '2025-01-31'
group by
  dimension_1,
  dimension_2,
  unblended_cost_unit
order by
  total_cost desc
limit 10;
```

### Cost trend analysis by meter category

Analyze cost trends by meter category and service to understand detailed billing patterns.

```sql+postgres
select
  usage_date,
  dimension_1 as meter_category,
  dimension_2 as service_name,
  unblended_cost_amount,
  lag(unblended_cost_amount) over (
    partition by dimension_1, dimension_2
    order by usage_date
  ) as previous_period_cost,
  unblended_cost_unit
from
  azure_cost_usage
where
  granularity = 'MONTHLY'
  and dimension_type_1 = 'MeterCategory'
  and dimension_type_2 = 'ServiceName'
order by
  dimension_1,
  dimension_2,
  usage_date desc;
```

```sql+sqlite
select
  usage_date,
  dimension_1 as meter_category,
  dimension_2 as service_name,
  unblended_cost_amount,
  lag(unblended_cost_amount) over (
    partition by dimension_1, dimension_2
    order by usage_date
  ) as previous_period_cost,
  unblended_cost_unit
from
  azure_cost_usage
where
  granularity = 'MONTHLY'
  and dimension_type_1 = 'MeterCategory'
  and dimension_type_2 = 'ServiceName'
order by
  dimension_1,
  dimension_2,
  usage_date desc;
```

### Cross-dimensional cost summary

Get a summary of costs across different dimension combinations to understand overall spending distribution.

```sql+postgres
select
  dimension_type_1,
  dimension_type_2,
  granularity,
  count(*) as total_records,
  sum(unblended_cost_amount) as total_cost,
  avg(unblended_cost_amount) as avg_cost,
  max(unblended_cost_amount) as max_cost,
  unblended_cost_unit
from
  azure_cost_usage
where
  usage_date >= current_date - interval '7 days'
group by
  dimension_type_1,
  dimension_type_2,
  granularity,
  unblended_cost_unit
order by
  total_cost desc;
```

```sql+sqlite
select
  dimension_type_1,
  dimension_type_2,
  granularity,
  count(*) as total_records,
  sum(unblended_cost_amount) as total_cost,
  avg(unblended_cost_amount) as avg_cost,
  max(unblended_cost_amount) as max_cost,
  unblended_cost_unit
from
  azure_cost_usage
where
  usage_date >= date('now', '-7 days')
group by
  dimension_type_1,
  dimension_type_2,
  granularity,
  unblended_cost_unit
order by
  total_cost desc;
```

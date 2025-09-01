---
title: "Steampipe Table: azure_cost_usage - Query Azure Cost and Usage Data using SQL"
description: "Allows users to query Azure Cost and Usage Data with flexible dimensions, providing detailed cost breakdown by any combination of Azure dimensions."
folder: "Cost Management"
---

# Table: azure_cost_usage - Query Azure Cost and Usage Data using SQL

Azure Cost Management provides cost analytics to help you understand and manage your Azure spending. The cost usage table allows you to query cost data with flexible dimensions, providing insights into Azure costs broken down by any combination of supported dimensions such as service name, resource group, location, resource type, subscription, and more. This enables custom cost analysis tailored to your organizational needs.

## Table Usage Guide

The `azure_cost_usage` table provides insights into cost and usage data within Microsoft Azure with flexible dimension support. As a Cloud Architect, FinOps engineer, or DevOps professional, explore cost details through this table using any combination of Azure dimensions. Utilize it to create custom cost breakdowns, analyze spending patterns across multiple dimensions, track costs by location and service, and perform advanced cost analytics that match your organizational structure.

**Important Notes:**
- You **_must_** specify `cost_type` (ActualCost or AmortizedCost), `granularity` (DAILY or MONTHLY), and at least one dimension qualifier — either `dimension_type_1` and/or `dimension_type_2` or `dimension_types` — in a `where` clause in order to use this table.
- For improved performance, it is advised that you use the optional quals `period_start` and `period_end` to limit the result set to a specific time period.
- This table supports optional quals. Queries with optional quals are optimised to use Azure Cost Management filters. Optional quals are supported for the following columns:
  - `scope` with supported operators `=`. Default to current subscription. Possible value are see: [Supported Scope](https://learn.microsoft.com/en-gb/rest/api/cost-management/query/usage?view=rest-cost-management-2025-03-01&tabs=HTTP#uri-parameters)
  - `period_start` with supported operators `=`. Default: 1 year ago.
  - `period_end` with supported operators `=`. Default: yesterday.

## Examples

### Recent daily costs by service and resource group
Get the last 7 days of daily costs broken down by service name and resource group to understand which services are costing the most in each resource group.

```sql+postgres
select
  usage_date,
  dimension_1 as service_name,
  dimension_2 as resource_group,
  cost,
  currency
from
  azure_cost_usage
where
  cost_type = 'ActualCost'
  and granularity = 'Daily'
  and dimension_type_1 = 'ServiceName'
  and dimension_type_2 = 'ResourceGroupName'
  and usage_date >= NOW() - INTERVAL '7 days'
order by
  usage_date desc,
  cost desc
limit 5;
```

```sql+sqlite
select
  usage_date,
  dimension_1 as service_name,
  dimension_2 as resource_group,
  cost,
  currency
from
  azure_cost_usage
where
  cost_type = 'ActualCost'
  and granularity = 'Daily'
  and dimension_type_1 = 'ServiceName'
  and dimension_type_2 = 'ResourceGroupName'
  and usage_date >= date('now', '-7 days')
order by
  usage_date desc,
  cost desc
limit 5;
```

### Location-based cost analysis (last 6 months)
Analyze the last 6 months of costs broken down by resource location to understand geographical spending patterns and identify the most expensive regions.

```sql+postgres
select
  dimension_1 as resource_location,
  sum(cost) as total_cost,
  currency
from
  azure_cost_usage
where
  cost_type = 'ActualCost'
  and granularity = 'Monthly'
  and dimension_type_1 = 'ResourceLocation'
  and usage_date >= NOW() - INTERVAL '6 months'
group by
  dimension_1,
  currency
order by
  total_cost desc
limit 10;
```

```sql+sqlite
select
  dimension_1 as resource_location,
  sum(cost) as total_cost,
  currency
from
  azure_cost_usage
where
  cost_type = 'ActualCost'
  and granularity = 'Monthly'
  and dimension_type_1 = 'ResourceLocation'
  and usage_date >= date('now', '-6 months')
group by
  dimension_1,
  currency
order by
  total_cost desc
limit 10;
```

### Costs for a specific billing period
Use period_start and period_end parameters to query costs for a specific monthly billing period (August 2025), showing service and resource group breakdown.

```sql+postgres
select
  usage_date,
  dimension_1 as service_name,
  dimension_2 as resource_group,
  cost,
  currency,
  period_start,
  period_end
from
  azure_cost_usage
where
  cost_type = 'ActualCost'
  and granularity = 'Daily'
  and dimension_type_1 = 'ServiceName'
  and dimension_type_2 = 'ResourceGroupName'
  and period_start = '2025-08-01'
  and period_end = '2025-08-31'
order by
  cost desc;
```

```sql+sqlite
select
  usage_date,
  dimension_1 as service_name,
  dimension_2 as resource_group,
  cost,
  currency,
  period_start,
  period_end
from
  azure_cost_usage
where
  cost_type = 'ActualCost'
  and granularity = 'Daily'
  and dimension_type_1 = 'ServiceName'
  and dimension_type_2 = 'ResourceGroupName'
  and period_start = '2025-08-01'
  and period_end = '2025-08-31'
order by
  cost desc;
```

### Daily service cost breakdown
Analyze the last 7 days of costs broken down by service name to understand which services are generating the highest costs.

```sql+postgres
select
  dimension_type_1,
  dimension_1,
  cost,
  currency,
  usage_date
from
  azure_cost_usage
where
  cost_type = 'ActualCost'
  and granularity = 'Daily'
  and dimension_type_1 = 'ServiceName'
  and usage_date >= NOW() - INTERVAL '7 days'
order by
  usage_date desc,
  cost desc
limit 10;
```

```sql+sqlite
select
  dimension_type_1,
  dimension_1,
  cost,
  currency,
  usage_date
from
  azure_cost_usage
where
  cost_type = 'ActualCost'
  and granularity = 'Daily'
  and dimension_type_1 = 'ServiceName'
  and usage_date >= date('now', '-7 days')
order by
  usage_date desc,
  cost desc
limit 10;
```

### Multi-subscription daily cost tracking
Track daily costs across different subscriptions and services over the last 30 days to understand multi-subscription spending patterns.

```sql+postgres
select
  usage_date,
  dimension_1 as subscription_name,
  dimension_2 as service_name,
  sum(cost) as total_cost,
  currency
from
  azure_cost_usage
where
  cost_type = 'ActualCost'
  and granularity = 'DAILY'
  and dimension_type_1 = 'SubscriptionName'
  and dimension_type_2 = 'ServiceName'
  and usage_date >= current_date - interval '30 days'
group by
  usage_date,
  dimension_1,
  dimension_2,
  currency
order by
  usage_date desc,
  total_cost desc;
```

```sql+sqlite
select
  usage_date,
  dimension_1 as subscription_name,
  dimension_2 as service_name,
  sum(cost) as total_cost,
  currency
from
  azure_cost_usage
where
  cost_type = 'ActualCost'
  and granularity = 'DAILY'
  and dimension_type_1 = 'SubscriptionName'
  and dimension_type_2 = 'ServiceName'
  and usage_date >= date('now', '-30 days')
group by
  usage_date,
  dimension_1,
  dimension_2,
  currency
order by
  usage_date desc,
  total_cost desc;
```

### Top service and resource group costs for January 2025
Find the highest cost service and resource group combinations within January 2025, showing total costs and number of usage days.

```sql+postgres
select
  dimension_1,
  dimension_2,
  sum(cost) as total_cost,
  count(*) as usage_days,
  currency
from
  azure_cost_usage
where
  cost_type = 'ActualCost'
  and granularity = 'DAILY'
  and dimension_type_1 = 'ServiceName'
  and dimension_type_2 = 'ResourceGroupName'
  and usage_date between '2025-01-01' and '2025-01-31'
group by
  dimension_1,
  dimension_2,
  currency
order by
  total_cost desc
limit 10;
```

```sql+sqlite
select
  dimension_1,
  dimension_2,
  sum(cost) as total_cost,
  count(*) as usage_days,
  currency
from
  azure_cost_usage
where
  cost_type = 'ActualCost'
  and granularity = 'DAILY'
  and dimension_type_1 = 'ServiceName'
  and dimension_type_2 = 'ResourceGroupName'
  and usage_date between '2025-01-01' and '2025-01-31'
group by
  dimension_1,
  dimension_2,
  currency
order by
  total_cost desc
limit 10;
```

### Monthly cost trends by meter category and service
Analyze monthly cost trends by meter category and service using window functions to compare current and previous period costs.

```sql+postgres
select
  usage_date,
  dimension_1 as meter_category,
  dimension_2 as service_name,
  cost,
  lag(cost) over (
    partition by dimension_1, dimension_2
    order by usage_date
  ) as previous_period_cost,
  currency
from
  azure_cost_usage
where
  cost_type = 'ActualCost'
  and granularity = 'MONTHLY'
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
  cost,
  lag(cost) over (
    partition by dimension_1, dimension_2
    order by usage_date
  ) as previous_period_cost,
  currency
from
  azure_cost_usage
where
  cost_type = 'ActualCost'
  and granularity = 'MONTHLY'
  and dimension_type_1 = 'MeterCategory'
  and dimension_type_2 = 'ServiceName'
order by
  dimension_1,
  dimension_2,
  usage_date desc;
```

### Weekly cost statistics by service and resource group
Get statistical summary (count, sum, average, max) of costs for service and resource group dimensions over the last 7 days.

```sql+postgres
select
  dimension_type_1,
  dimension_type_2,
  granularity,
  count(*) as total_records,
  sum(cost) as total_cost,
  avg(cost) as avg_cost,
  max(cost) as max_cost,
  currency
from
  azure_cost_usage
where
  cost_type = 'ActualCost'
  and granularity = 'DAILY'
  and dimension_type_1 = 'ServiceName'
  and dimension_type_2 = 'ResourceGroupName'
  and usage_date >= current_date - interval '7 days'
group by
  dimension_type_1,
  dimension_type_2,
  granularity,
  currency
order by
  total_cost desc;
```

```sql+sqlite
select
  dimension_type_1,
  dimension_type_2,
  granularity,
  count(*) as total_records,
  sum(cost) as total_cost,
  avg(cost) as avg_cost,
  max(cost) as max_cost,
  currency
from
  azure_cost_usage
where
  cost_type = 'ActualCost'
  and granularity = 'DAILY'
  and dimension_type_1 = 'ServiceName'
  and dimension_type_2 = 'ResourceGroupName'
  and usage_date >= date('now', '-7 days')
group by
  dimension_type_1,
  dimension_type_2,
  granularity,
  currency
order by
  total_cost desc;
```

### Reservation savings analysis
Compare actual costs vs amortized costs across service and resource group dimensions to identify reservation savings and understand cost optimization opportunities.

```sql+postgres
with actual_costs as (
  select
    dimension_1 as service_name,
    dimension_2 as resource_group,
    usage_date,
    cost as actual_cost,
    currency
  from
    azure_cost_usage
  where
    cost_type = 'ActualCost'
    and granularity = 'DAILY'
    and dimension_type_1 = 'ServiceName'
    and dimension_type_2 = 'ResourceGroupName'
),
amortized_costs as (
  select
    dimension_1 as service_name,
    dimension_2 as resource_group,
    usage_date,
    cost as amortized_cost,
    currency
  from
    azure_cost_usage
  where
    cost_type = 'AmortizedCost'
    and granularity = 'DAILY'
    and dimension_type_1 = 'ServiceName'
    and dimension_type_2 = 'ResourceGroupName'
)
select
  a.service_name,
  a.resource_group,
  a.usage_date,
  a.actual_cost,
  am.amortized_cost,
  (a.actual_cost - am.amortized_cost) as reservation_savings,
  a.currency
from
  actual_costs a
  join amortized_costs am on a.service_name = am.service_name
  and a.resource_group = am.resource_group
  and a.usage_date = am.usage_date
  and a.currency = am.currency
where
  a.actual_cost != am.amortized_cost
order by
  reservation_savings desc;
```

```sql+sqlite
with actual_costs as (
  select
    dimension_1 as service_name,
    dimension_2 as resource_group,
    usage_date,
    cost as actual_cost,
    currency
  from
    azure_cost_usage
  where
    cost_type = 'ActualCost'
    and granularity = 'DAILY'
    and dimension_type_1 = 'ServiceName'
    and dimension_type_2 = 'ResourceGroupName'
),
amortized_costs as (
  select
    dimension_1 as service_name,
    dimension_2 as resource_group,
    usage_date,
    cost as amortized_cost,
    currency
  from
    azure_cost_usage
  where
    cost_type = 'AmortizedCost'
    and granularity = 'DAILY'
    and dimension_type_1 = 'ServiceName'
    and dimension_type_2 = 'ResourceGroupName'
)
select
  a.service_name,
  a.resource_group,
  a.usage_date,
  a.actual_cost,
  am.amortized_cost,
  (a.actual_cost - am.amortized_cost) as reservation_savings,
  a.currency
from
  actual_costs a
  join amortized_costs am on a.service_name = am.service_name
  and a.resource_group = am.resource_group
  and a.usage_date = am.usage_date
  and a.currency = am.currency
where
  a.actual_cost != am.amortized_cost
order by
  reservation_savings desc;
```

### Monthly amortized cost across multiple dimensions
Show monthly aggregated amortized cost with a multi-dimension breakdown; returns cost, usage_date (month-end), and a pretty-printed dimensions JSON.

```sql+postgres
select
  cost,
  usage_date,
  jsonb_pretty(dimensions)
from
  azure_cost_usage
where
  granularity = 'MONTHLY'
  and cost_type = 'AmortizedCost'
  and dimension_types = '["ResourceGroupName", "ServiceName", "ResourceLocation", "ResourceId", "MeterCategory", "ResourceType", "ChargeType"]';
```

```sql+sqlite
select
  cost,
  usage_date,
  json_pretty(dimensions) as dimensions
from
  azure_cost_usage
where
  granularity = 'MONTHLY'
  and cost_type = 'AmortizedCost'
  and dimension_types = '["ResourceGroupName", "ServiceName", "ResourceLocation", "ResourceId", "MeterCategory", "ResourceType", "ChargeType"]';
```

---
title: "Steampipe Table: azure_cost_by_service_daily - Query Azure Daily Service Costs using SQL"
description: "Allows users to query Azure Daily Service Costs, providing detailed cost breakdown by service on a daily basis."
folder: "Cost Management"
---

# Table: azure_cost_by_service_daily - Query Azure Daily Service Costs using SQL

Azure Cost Management provides cost analytics to help you understand and manage your Azure spending. The daily service cost breakdown allows you to track costs at a granular level, providing insights into how much each Azure service costs on a day-by-day basis. This helps in identifying cost trends, optimizing resource usage, and managing budgets effectively.

## Table Usage Guide

The `azure_cost_by_service_daily` table provides insights into daily cost breakdown by service within Microsoft Azure. As a Cloud Architect, FinOps engineer, or DevOps professional, explore service-specific cost details through this table, including daily usage costs, currency information, and service names. Utilize it to uncover cost patterns, identify expensive services, track daily spending trends, and optimize resource allocation.

**Important Notes:**

- You **_must_** specify `cost_type` (ActualCost or AmortizedCost) in a `where` clause in order to use this table.
- For improved performance, it is advised that you use the optional quals `period_start` and `period_end` to limit the result set to a specific time period.
- This table supports optional quals. Queries with optional quals are optimised to use Azure Cost Management filters. Optional quals are supported for the following columns:
  - `scope` with supported operators `=`. Default to current subscription. Possible value are see: [Supported Scope](https://learn.microsoft.com/en-gb/rest/api/cost-management/query/usage?view=rest-cost-management-2025-03-01&tabs=HTTP#uri-parameters)
  - `period_start` with supported operators `=`. Default: 1 year ago.
  - `period_end` with supported operators `=`. Default: yesterday.
  - `service_name` with supported operators `=`, `<>`.

## Examples

### Recent daily costs by service
Get the most recent 7 days of daily costs across Azure services, showing the cost breakdown by service with currency details.

```sql+postgres
select
  service_name,
  usage_date,
  cost,
  currency
from
  azure_cost_by_service_daily
where
  cost_type = 'ActualCost'
  and usage_date >= NOW() - INTERVAL '7 days'
order by
  usage_date desc,
  cost desc
limit 10;
```

```sql+sqlite
select
  service_name,
  usage_date,
  cost,
  currency
from
  azure_cost_by_service_daily
where
  cost_type = 'ActualCost'
  and usage_date >= date('now', '-7 days')
order by
  usage_date desc,
  cost desc
limit 10;
```

### Historical daily costs for a specific service
Analyze the complete historical daily cost trend for a specific Azure service to understand its usage patterns and cost evolution over time.

```sql+postgres
select
  usage_date,
  service_name,
  cost,
  pre_tax_cost,
  currency
from
  azure_cost_by_service_daily
where
  cost_type = 'ActualCost'
  and service_name = 'Microsoft Defender for Cloud'
order by
  usage_date desc;
```

```sql+sqlite
select
  usage_date,
  service_name,
  cost,
  pre_tax_cost,
  currency
from
  azure_cost_by_service_daily
where
  cost_type = 'ActualCost'
  and service_name = 'Microsoft Defender for Cloud'
order by
  usage_date desc;
```

### Costs for a specific billing period
Use period_start and period_end parameters to query costs for a specific time range, showing both actual cost and pre-tax cost with period metadata.

```sql+postgres
select
  usage_date,
  service_name,
  cost,
  pre_tax_cost,
  currency,
  period_start,
  period_end
from
  azure_cost_by_service_daily
where
  cost_type = 'ActualCost'
  and period_start = '2025-08-01'
  and period_end = '2025-08-31'
order by
  cost desc;
```

```sql+sqlite
select
  usage_date,
  service_name,
  cost,
  pre_tax_cost,
  currency,
  period_start,
  period_end
from
  azure_cost_by_service_daily
where
  cost_type = 'ActualCost'
  and period_start = '2025-08-01'
  and period_end = '2025-08-31'
order by
  cost desc;
```

### Service cost aggregation analysis
Identify the top 5 most expensive Azure services over the last 30 days with total costs and average daily spending to focus optimization efforts.

```sql+postgres
select
  service_name,
  sum(cost) as total_cost,
  avg(cost) as avg_daily_cost,
  currency
from
  azure_cost_by_service_daily
where
  cost_type = 'ActualCost'
  and usage_date >= NOW() - INTERVAL '30 days'
group by
  service_name,
  currency
order by
  total_cost desc
limit 5;
```

```sql+sqlite
select
  service_name,
  sum(cost) as total_cost,
  avg(cost) as avg_daily_cost,
  currency
from
  azure_cost_by_service_daily
where
  cost_type = 'ActualCost'
  and usage_date >= date('now', '-30 days')
group by
  service_name,
  currency
order by
  total_cost desc
limit 5;
```

### Daily total spending trends
Analyze the aggregated daily cost trends across all services to identify spending patterns and cost spikes over the last 30 days.

```sql+postgres
select
  usage_date,
  sum(cost) as daily_total_cost,
  currency
from
  azure_cost_by_service_daily
where
  cost_type = 'ActualCost'
  and usage_date >= NOW() - INTERVAL '30 days'
group by
  usage_date,
  currency
order by
  usage_date desc;
```

```sql+sqlite
select
  usage_date,
  sum(cost) as daily_total_cost,
  currency
from
  azure_cost_by_service_daily
where
  cost_type = 'ActualCost'
  and usage_date >= date('now', '-30 days')
group by
  usage_date,
  currency
order by
  usage_date desc;
```

### Service costs within a date range
Get detailed daily cost breakdown for all services within a specific 7-day period, ordered by usage date and cost amount.

```sql+postgres
select
  usage_date,
  service_name,
  cost,
  currency
from
  azure_cost_by_service_daily
where
  cost_type = 'ActualCost'
  and usage_date between '2025-08-25' and '2025-08-31'
order by
  usage_date,
  cost desc;
```

```sql+sqlite
select
  usage_date,
  service_name,
  cost,
  currency
from
  azure_cost_by_service_daily
where
  cost_type = 'ActualCost'
  and usage_date between '2025-08-25' and '2025-08-31'
order by
  usage_date,
  cost desc;
```

### High-cost services above threshold
Identify services that exceeded a specific cost threshold ($0.50) on any given day, useful for cost monitoring and budget alerting.

```sql+postgres
select
  usage_date,
  service_name,
  cost,
  currency
from
  azure_cost_by_service_daily
where
  cost_type = 'ActualCost'
  and cost > 0.5
order by
  cost desc;
```

```sql+sqlite
select
  usage_date,
  service_name,
  cost,
  currency
from
  azure_cost_by_service_daily
where
  cost_type = 'ActualCost'
  and cost > 0.5
order by
  cost desc;
```

### Reservation savings analysis
Compare actual costs vs amortized costs to identify reservation savings by joining ActualCost and AmortizedCost data for each service and day.

```sql+postgres
with actual_costs as (
  select
    service_name,
    usage_date,
    cost as actual_cost,
    currency
  from
    azure_cost_by_service_daily
  where
    cost_type = 'ActualCost'
),
amortized_costs as (
  select
    service_name,
    usage_date,
    cost as amortized_cost,
    currency
  from
    azure_cost_by_service_daily
  where
    cost_type = 'AmortizedCost'
)
select
  a.service_name,
  a.usage_date,
  a.actual_cost,
  am.amortized_cost,
  (a.actual_cost - am.amortized_cost) as reservation_savings,
  a.currency
from
  actual_costs a
  join amortized_costs am on a.service_name = am.service_name
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
    service_name,
    usage_date,
    cost as actual_cost,
    currency
  from
    azure_cost_by_service_daily
  where
    cost_type = 'ActualCost'
),
amortized_costs as (
  select
    service_name,
    usage_date,
    cost as amortized_cost,
    currency
  from
    azure_cost_by_service_daily
  where
    cost_type = 'AmortizedCost'
)
select
  a.service_name,
  a.usage_date,
  a.actual_cost,
  am.amortized_cost,
  (a.actual_cost - am.amortized_cost) as reservation_savings,
  a.currency
from
  actual_costs a
  join amortized_costs am on a.service_name = am.service_name
  and a.usage_date = am.usage_date
  and a.currency = am.currency
where
  a.actual_cost != am.amortized_cost
order by
  reservation_savings desc;
```

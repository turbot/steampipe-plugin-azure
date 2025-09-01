---
title: "Steampipe Table: azure_cost_by_service_monthly - Query Azure Monthly Service Costs using SQL"
description: "Allows users to query Azure Monthly Service Costs, providing detailed cost breakdown by service on a monthly basis."
folder: "Cost Management"
---

# Table: azure_cost_by_service_monthly - Query Azure Monthly Service Costs using SQL

Azure Cost Management provides cost analytics to help you understand and manage your Azure spending. The monthly service cost breakdown allows you to track costs at a higher level, providing insights into how much each Azure service costs on a month-by-month basis. This helps in understanding long-term cost trends, budget planning, and strategic cost management decisions.

## Table Usage Guide

The `azure_cost_by_service_monthly` table provides insights into monthly cost breakdown by service within Microsoft Azure. As a Cloud Architect, FinOps engineer, or DevOps professional, explore service-specific monthly cost details through this table, including monthly usage costs, currency information, and service names. Utilize it to uncover monthly cost patterns, compare service costs across months, track long-term spending trends, and support budget planning and forecasting.

**Important Notes:**

- You **_must_** specify `cost_type` (ActualCost or AmortizedCost) in a `where` clause in order to use this table.
- For improved performance, it is advised that you use the optional quals `period_start` and `period_end` to limit the result set to a specific time period.
- This table supports optional quals. Queries with optional quals are optimised to use Azure Cost Management filters. Optional quals are supported for the following columns:
  - `scope` with supported operators `=`. Default to current subscription. Possible value are see: [Supported Scope](https://learn.microsoft.com/en-gb/rest/api/cost-management/query/usage?view=rest-cost-management-2025-03-01&tabs=HTTP#uri-parameters)
  - `period_start` with supported operators `=`. Default: 1 year ago.
  - `period_end` with supported operators `=`. Default: yesterday.
  - `service_name` with supported operators `=`, `<>`.

## Examples

### Recent monthly costs by service

Get the last 6 months of monthly costs across Azure services, showing the cost breakdown by service with currency details.

```sql+postgres
select
  service_name,
  usage_date,
  cost,
  currency
from
  azure_cost_by_service_monthly
where
  cost_type = 'ActualCost'
  and usage_date >= NOW() - INTERVAL '6 months'
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
  azure_cost_by_service_monthly
where
  cost_type = 'ActualCost'
  and usage_date >= date('now', '-6 months')
order by
  usage_date desc,
  cost desc
limit 10;
```

### Historical monthly costs for a specific service

Analyze the complete historical monthly cost trend for a specific Azure service to understand its usage patterns and cost evolution over time.

```sql+postgres
select
  usage_date,
  service_name,
  cost,
  pre_tax_cost,
  currency
from
  azure_cost_by_service_monthly
where
  cost_type = 'ActualCost'
  and service_name = 'Storage'
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
  azure_cost_by_service_monthly
where
  cost_type = 'ActualCost'
  and service_name = 'Storage'
order by
  usage_date desc;
```

### Costs for a specific billing period

Use period_start and period_end parameters to query costs for a specific annual billing period (2024), showing all services and their costs.

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
  azure_cost_by_service_monthly
where
  cost_type = 'ActualCost'
  and period_start = '2024-01-01'
  and period_end = '2024-12-31'
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
  azure_cost_by_service_monthly
where
  cost_type = 'ActualCost'
  and period_start = '2024-01-01'
  and period_end = '2024-12-31'
order by
  cost desc;
```

### Total historical spend by service

Get the total cumulative cost for each service across all months to understand which services contribute most to your overall Azure spending.

```sql+postgres
select
  service_name,
  sum(cost) as total_monthly_cost,
  currency,
  count(*) as months_with_usage
from
  azure_cost_by_service_monthly
where
  cost_type = 'ActualCost'
group by
  service_name,
  currency
order by
  total_monthly_cost desc;
```

```sql+sqlite
select
  service_name,
  sum(cost) as total_monthly_cost,
  currency,
  count(*) as months_with_usage
from
  azure_cost_by_service_monthly
where
  cost_type = 'ActualCost'
group by
  service_name,
  currency
order by
  total_monthly_cost desc;
```

### Monthly cost trends by year

Compare monthly costs between 2024 and 2025 to understand year-over-year cost growth and seasonal spending patterns.

```sql+postgres
select
  extract(month from usage_date) as month_number,
  extract(year from usage_date) as year,
  sum(cost) as monthly_total,
  currency
from
  azure_cost_by_service_monthly
where
  cost_type = 'ActualCost'
  and extract(year from usage_date) in (2024, 2025)
group by
  extract(month from usage_date),
  extract(year from usage_date),
  currency
order by
  month_number,
  year;
```

```sql+sqlite
select
  cast(strftime('%m', usage_date) as integer) as month_number,
  cast(strftime('%Y', usage_date) as integer) as year,
  sum(cost) as monthly_total,
  currency
from
  azure_cost_by_service_monthly
where
  cost_type = 'ActualCost'
  and cast(strftime('%Y', usage_date) as integer) in (2024, 2025)
group by
  cast(strftime('%m', usage_date) as integer),
  cast(strftime('%Y', usage_date) as integer),
  currency
order by
  month_number,
  year;
```

### Services with month-over-month cost increases

Identify services where costs increased from the previous month, showing the actual cost change amounts to help focus cost optimization efforts.

```sql+postgres
with monthly_costs as (
  select
    service_name,
    usage_date,
    cost,
    lag(cost) over (partition by service_name order by usage_date) as prev_month_cost
  from
    azure_cost_by_service_monthly
  where
    cost_type = 'ActualCost'
)
select
  service_name,
  usage_date,
  cost as current_month,
  prev_month_cost as previous_month,
  cost - prev_month_cost as cost_change
from
  monthly_costs
where
  prev_month_cost is not null
  and cost > prev_month_cost
order by
  cost_change desc;
```

```sql+sqlite
with monthly_costs as (
  select
    service_name,
    usage_date,
    cost,
    lag(cost) over (partition by service_name order by usage_date) as prev_month_cost
  from
    azure_cost_by_service_monthly
  where
    cost_type = 'ActualCost'
)
select
  service_name,
  usage_date,
  cost as current_month,
  prev_month_cost as previous_month,
  cost - prev_month_cost as cost_change
from
  monthly_costs
where
  prev_month_cost is not null
  and cost > prev_month_cost
order by
  cost_change desc;
```

### Cost statistics by service

Calculate average, minimum, and maximum monthly costs for each service to understand spending patterns and cost variability.

```sql+postgres
select
  service_name,
  avg(cost) as avg_monthly_cost,
  min(cost) as min_monthly_cost,
  max(cost) as max_monthly_cost,
  currency
from
  azure_cost_by_service_monthly
where
  cost_type = 'ActualCost'
group by
  service_name,
  currency
order by
  avg_monthly_cost desc;
```

```sql+sqlite
select
  service_name,
  round(avg(cost), 2) as avg_monthly_cost,
  round(min(cost), 2) as min_monthly_cost,
  round(max(cost), 2) as max_monthly_cost,
  currency
from
  azure_cost_by_service_monthly
where
  cost_type = 'ActualCost'
group by
  service_name,
  currency
order by
  avg_monthly_cost desc;
```

### Tax analysis for amortized costs

Analyze the difference between pre-tax costs and final amortized costs to understand tax implications on monthly service spending.

```sql+postgres
select
  service_name,
  usage_date,
  pre_tax_cost,
  cost,
  (pre_tax_cost - cost) as tax_amount,
  currency
from
  azure_cost_by_service_monthly
where
  cost_type = 'AmortizedCost'
  and pre_tax_cost is not null
  and pre_tax_cost != cost
order by
  tax_amount desc;
```

```sql+sqlite
select
  service_name,
  usage_date,
  pre_tax_cost,
  cost,
  (pre_tax_cost - cost) as tax_amount,
  currency
from
  azure_cost_by_service_monthly
where
  cost_type = 'AmortizedCost'
  and pre_tax_cost is not null
  and pre_tax_cost != cost
order by
  tax_amount desc;
```

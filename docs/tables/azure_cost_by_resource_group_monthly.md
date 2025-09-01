---
title: "Steampipe Table: azure_cost_by_resource_group_monthly - Query Azure Monthly Resource Group Costs using SQL"
description: "Allows users to query Azure Monthly Resource Group Costs, providing detailed cost breakdown by resource group on a monthly basis."
folder: "Cost Management"
---

# Table: azure_cost_by_resource_group_monthly - Query Azure Monthly Resource Group Costs using SQL

Azure Cost Management provides cost analytics to help you understand and manage your Azure spending. The monthly resource group cost breakdown allows you to track costs at a higher level, providing insights into how much each Azure resource group costs on a month-by-month basis. This helps in understanding team/project allocations, long-term trends, and budget planning.

## Table Usage Guide

The `azure_cost_by_resource_group_monthly` table provides insights into monthly cost breakdown by resource group within Microsoft Azure. As a Cloud Architect, FinOps engineer, or DevOps professional, explore resource group-specific monthly cost details through this table, including monthly usage costs, currency information, and resource group names. Utilize it to uncover monthly cost patterns and support budget planning and forecasting.

**Important Notes:**

- You **_must_** specify `cost_type` (ActualCost or AmortizedCost) in a `where` clause in order to use this table.
- For improved performance, it is advised that you use the optional quals `period_start` and `period_end` to limit the result set to a specific time period.
- This table supports optional quals. Queries with optional quals are optimised to use Azure Cost Management filters. Optional quals are supported for the following columns:
  - `scope` with supported operators `=`. Default to current subscription. Possible value are see: [Supported Scope](https://learn.microsoft.com/en-gb/rest/api/cost-management/query/usage?view=rest-cost-management-2025-03-01&tabs=HTTP#uri-parameters)
  - `period_start` with supported operators `=`. Default: 1 year ago.
  - `period_end` with supported operators `=`. Default: yesterday.
  - `resource_group` with supported operators `=`, `<>`.

## Examples

### Recent monthly costs by resource group
Get the last 6 months of monthly costs across Azure resource groups, showing the cost breakdown by resource group with subscription details.

```sql+postgres
select
  resource_group,
  usage_date,
  cost,
  currency
from
  azure_cost_by_resource_group_monthly
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
  resource_group,
  usage_date,
  cost,
  currency
from
  azure_cost_by_resource_group_monthly
where
  cost_type = 'ActualCost'
  and usage_date >= date('now', '-6 months')
order by
  usage_date desc,
  cost desc
limit 10;
```

### Historical monthly costs for a specific resource group
Analyze the complete historical monthly cost trend for a specific Azure resource group to understand its usage patterns and cost evolution over time.

```sql+postgres
select
  usage_date,
  resource_group,
  cost,
  pre_tax_cost,
  currency
from
  azure_cost_by_resource_group_monthly
where
  cost_type = 'ActualCost'
  and resource_group = 'demo'
order by
  usage_date desc;
```

```sql+sqlite
select
  usage_date,
  resource_group,
  cost,
  pre_tax_cost,
  currency
from
  azure_cost_by_resource_group_monthly
where
  cost_type = 'ActualCost'
  and resource_group = 'demo'
order by
  usage_date desc;
```

### Costs for a specific billing period
Use period_start and period_end parameters to query costs for a specific time range, showing both actual cost and pre-tax cost with period metadata.

```sql+postgres
select
  usage_date,
  resource_group,
  cost,
  pre_tax_cost,
  currency,
  period_start,
  period_end
from
  azure_cost_by_resource_group_monthly
where
  cost_type = 'ActualCost'
  and period_start = '2025-01-01'
  and period_end = '2025-12-31'
order by
  cost desc;
```

```sql+sqlite
select
  usage_date,
  resource_group,
  cost,
  pre_tax_cost,
  currency,
  period_start,
  period_end
from
  azure_cost_by_resource_group_monthly
where
  cost_type = 'ActualCost'
  and period_start = '2025-01-01'
  and period_end = '2025-12-31'
order by
  cost desc;
```

### Aggregate spending analysis by resource group
Get the total cumulative costs across all months for each resource group to identify which teams or projects contribute most to your overall Azure spending.

```sql+postgres
select
  resource_group,
  sum(cost) as total_monthly_cost,
  currency,
  count(*) as months_with_usage
from
  azure_cost_by_resource_group_monthly
where
  cost_type = 'ActualCost'
group by
  resource_group,
  currency
order by
  total_monthly_cost desc
limit 10;
```

```sql+sqlite
select
  resource_group,
  sum(cost) as total_monthly_cost,
  currency,
  count(*) as months_with_usage
from
  azure_cost_by_resource_group_monthly
where
  cost_type = 'ActualCost'
group by
  resource_group,
  currency
order by
  total_monthly_cost desc
limit 10;
```

### Month-over-month cost increases
Identify resource groups with increasing costs by comparing each month to the previous month, showing absolute cost changes to help focus optimization efforts.

```sql+postgres
with monthly_costs as (
  select
    resource_group,
    usage_date,
    cost,
    lag(cost) over (partition by resource_group order by usage_date) as prev_month_cost
  from
    azure_cost_by_resource_group_monthly
  where
    cost_type = 'ActualCost'
)
select
  resource_group,
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
  cost_change desc
limit 10;
```

```sql+sqlite
with monthly_costs as (
  select
    resource_group,
    usage_date,
    cost,
    lag(cost) over (partition by resource_group order by usage_date) as prev_month_cost
  from
    azure_cost_by_resource_group_monthly
  where
    cost_type = 'ActualCost'
)
select
  resource_group,
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
  cost_change desc
limit 10;
```

### Budget threshold monitoring
Track monthly spending against predefined budget thresholds ($50 over budget, $30 near budget) by resource group to monitor departmental spending limits.

```sql+postgres
select
  resource_group,
  usage_date,
  cost,
  case
    when cost > 50 then 'Over Budget'
    when cost > 30 then 'Near Budget'
    else 'Within Budget'
  end as budget_status,
  currency
from
  azure_cost_by_resource_group_monthly
where
  cost_type = 'ActualCost'
order by
  usage_date desc,
  cost desc;
```

```sql+sqlite
select
  resource_group,
  usage_date,
  cost,
  case
    when cost > 50 then 'Over Budget'
    when cost > 30 then 'Near Budget'
    else 'Within Budget'
  end as budget_status,
  currency
from
  azure_cost_by_resource_group_monthly
where
  cost_type = 'ActualCost'
order by
  usage_date desc,
  cost desc;
```

### Cost statistics by resource group
Analyze cost distribution statistics (average, min, max, standard deviation) across resource groups to understand spending patterns and variability.

```sql+postgres
select
  resource_group,
  avg(cost) as avg_monthly_cost,
  min(cost) as min_monthly_cost,
  max(cost) as max_monthly_cost,
  stddev(cost) as cost_variability,
  currency
from
  azure_cost_by_resource_group_monthly
where
  cost_type = 'ActualCost'
group by
  resource_group,
  currency
order by
  avg_monthly_cost desc
limit 10;
```

```sql+sqlite
select
  resource_group,
  avg(cost) as avg_monthly_cost,
  min(cost) as min_monthly_cost,
  max(cost) as max_monthly_cost,
  currency
from
  azure_cost_by_resource_group_monthly
where
  cost_type = 'ActualCost'
group by
  resource_group,
  currency
order by
  avg_monthly_cost desc
limit 10;
```

### Quarterly cost aggregation
Aggregate monthly costs into quarterly totals for each resource group to analyze seasonal spending patterns and quarterly budget planning.

```sql+postgres
select
  resource_group,
  extract(year from usage_date) as year,
  extract(quarter from usage_date) as quarter,
  sum(cost) as quarterly_cost,
  currency
from
  azure_cost_by_resource_group_monthly
where
  cost_type = 'ActualCost'
group by
  resource_group,
  extract(year from usage_date),
  extract(quarter from usage_date),
  currency
order by
  resource_group,
  year,
  quarter;
```

```sql+sqlite
select
  resource_group,
  cast(strftime('%Y', usage_date) as integer) as year,
  case
    when cast(strftime('%m', usage_date) as integer) <= 3 then 1
    when cast(strftime('%m', usage_date) as integer) <= 6 then 2
    when cast(strftime('%m', usage_date) as integer) <= 9 then 3
    else 4
  end as quarter,
  sum(cost) as quarterly_cost,
  currency
from
  azure_cost_by_resource_group_monthly
where
  cost_type = 'ActualCost'
group by
  resource_group,
  cast(strftime('%Y', usage_date) as integer),
  case
    when cast(strftime('%m', usage_date) as integer) <= 3 then 1
    when cast(strftime('%m', usage_date) as integer) <= 6 then 2
    when cast(strftime('%m', usage_date) as integer) <= 9 then 3
    else 4
  end,
  currency
order by
  resource_group,
  year,
  quarter;
```

### Reservation savings analysis
Compare actual costs vs amortized costs to identify reservation savings by joining ActualCost and AmortizedCost data for each resource group and month.

```sql+postgres
with actual_costs as (
  select
    resource_group,
    usage_date,
    cost as actual_cost,
    currency
  from
    azure_cost_by_resource_group_monthly
  where
    cost_type = 'ActualCost'
),
amortized_costs as (
  select
    resource_group,
    usage_date,
    cost as amortized_cost,
    currency
  from
    azure_cost_by_resource_group_monthly
  where
    cost_type = 'AmortizedCost'
)
select
  a.resource_group,
  a.usage_date,
  a.actual_cost,
  am.amortized_cost,
  (a.actual_cost - am.amortized_cost) as reservation_savings,
  a.currency
from
  actual_costs a
  join amortized_costs am on a.resource_group = am.resource_group
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
    resource_group,
    usage_date,
    cost as actual_cost,
    currency
  from
    azure_cost_by_resource_group_monthly
  where
    cost_type = 'ActualCost'
),
amortized_costs as (
  select
    resource_group,
    usage_date,
    cost as amortized_cost,
    currency
  from
    azure_cost_by_resource_group_monthly
  where
    cost_type = 'AmortizedCost'
)
select
  a.resource_group,
  a.usage_date,
  a.actual_cost,
  am.amortized_cost,
  (a.actual_cost - am.amortized_cost) as reservation_savings,
  a.currency
from
  actual_costs a
  join amortized_costs am on a.resource_group = am.resource_group
  and a.usage_date = am.usage_date
  and a.currency = am.currency
where
  a.actual_cost != am.amortized_cost
order by
  reservation_savings desc;
```

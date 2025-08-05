---
title: "Steampipe Table: azure_cost_by_resource_group_monthly - Query Azure Monthly Resource Group Costs using SQL"
description: "Allows users to query Azure Monthly Resource Group Costs, providing detailed cost breakdown by resource group on a monthly basis."
folder: "Cost Management"
---

# Table: azure_cost_by_resource_group_monthly - Query Azure Monthly Resource Group Costs using SQL

Azure Cost Management provides cost analytics to help you understand and manage your Azure spending. The monthly resource group cost breakdown allows you to track costs at a higher level, providing insights into how much each Azure resource group costs on a month-by-month basis. This helps in understanding long-term cost trends per team or project, budget planning, and strategic cost management decisions for different departments.

## Table Usage Guide

The `azure_cost_by_resource_group_monthly` table provides insights into monthly cost breakdown by resource group within Microsoft Azure. As a Cloud Architect, FinOps engineer, or DevOps professional, explore resource group-specific monthly cost details through this table, including monthly usage costs, currency information, and resource group names. Utilize it to uncover monthly cost patterns per team or project, compare resource group costs across months, track long-term spending trends, and support departmental budget planning and forecasting.

**Important Notes:**

- This table supports optional quals. Queries with optional quals are optimised to reduce query time and improve performance. Optional quals are supported for the following columns:
  - `scope` with supported operators `=`.
  - `type` with supported operators `=`. Valid values are 'ActualCost' (default) and 'AmortizedCost'.
  - `period_start` with supported operators `=`, `>=`, `>`, `<=`, and `<`.
  - `period_end` with supported operators `=`, `>=`, `>`, `<=`, and `<`.

## Examples

### Basic monthly cost info by resource group

Explore monthly costs across different Azure resource groups to understand your spending patterns and identify the most expensive resource groups by month.

```sql+postgres
select
  usage_date,
  resource_group,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_monthly
order by
  usage_date desc,
  pre_tax_cost_amount desc;
```

```sql+sqlite
select
  usage_date,
  resource_group,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_monthly
order by
  usage_date desc,
  pre_tax_cost_amount desc;
```

### Monthly cost trend for a specific resource group

Analyze the monthly cost trend for a specific Azure resource group to understand its long-term usage patterns and cost evolution.

```sql+postgres
select
  usage_date,
  resource_group,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_monthly
where
  resource_group = 'production-rg'
order by
  usage_date desc;
```

```sql+sqlite
select
  usage_date,
  resource_group,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_monthly
where
  resource_group = 'production-rg'
order by
  usage_date desc;
```

### Query costs for a specific period

Use period_start and period_end parameters to query costs for a specific time range.

```sql+postgres
select
  usage_date,
  resource_group,
  pre_tax_cost_amount,
  pre_tax_cost_unit,
  period_start,
  period_end
from
  azure_cost_by_resource_group_monthly
where
  period_start = '2024-01-01'
  and period_end = '2024-12-31'
order by
  pre_tax_cost_amount desc;
```

```sql+sqlite
select
  usage_date,
  resource_group,
  pre_tax_cost_amount,
  pre_tax_cost_unit,
  period_start,
  period_end
from
  azure_cost_by_resource_group_monthly
where
  period_start = '2024-01-01'
  and period_end = '2024-12-31'
order by
  pre_tax_cost_amount desc;
```

### Total monthly spend by resource group

Get the total monthly cost for each resource group to understand which teams or projects contribute most to your monthly Azure bill.

```sql+postgres
select
  resource_group,
  sum(pre_tax_cost_amount) as total_monthly_cost,
  pre_tax_cost_unit,
  count(*) as months_with_usage
from
  azure_cost_by_resource_group_monthly
group by
  resource_group,
  pre_tax_cost_unit
order by
  total_monthly_cost desc;
```

```sql+sqlite
select
  resource_group,
  sum(pre_tax_cost_amount) as total_monthly_cost,
  pre_tax_cost_unit,
  count(*) as months_with_usage
from
  azure_cost_by_resource_group_monthly
group by
  resource_group,
  pre_tax_cost_unit
order by
  total_monthly_cost desc;
```

### Resource groups with increasing monthly costs

Identify resource groups where costs are trending upward by comparing recent months to help focus cost optimization efforts on specific teams or projects.

```sql+postgres
with monthly_costs as (
  select
    resource_group,
    usage_date,
    pre_tax_cost_amount,
    lag(pre_tax_cost_amount) over (partition by resource_group order by usage_date) as prev_month_cost
  from
    azure_cost_by_resource_group_monthly
)
select
  resource_group,
  usage_date,
  pre_tax_cost_amount as current_month,
  prev_month_cost as previous_month,
  pre_tax_cost_amount - prev_month_cost as cost_change
from
  monthly_costs
where
  prev_month_cost is not null
  and pre_tax_cost_amount > prev_month_cost
order by
  cost_change desc;
```

```sql+sqlite
with monthly_costs as (
  select
    resource_group,
    usage_date,
    pre_tax_cost_amount,
    lag(pre_tax_cost_amount) over (partition by resource_group order by usage_date) as prev_month_cost
  from
    azure_cost_by_resource_group_monthly
)
select
  resource_group,
  usage_date,
  pre_tax_cost_amount as current_month,
  prev_month_cost as previous_month,
  pre_tax_cost_amount - prev_month_cost as cost_change
from
  monthly_costs
where
  prev_month_cost is not null
  and pre_tax_cost_amount > prev_month_cost
order by
  cost_change desc;
```

### Departmental budget tracking

Track monthly spending against potential budgets by resource group to monitor departmental or team spending limits.

```sql+postgres
select
  resource_group,
  usage_date,
  pre_tax_cost_amount,
  case
    when pre_tax_cost_amount > 1000 then 'Over Budget'
    when pre_tax_cost_amount > 800 then 'Near Budget'
    else 'Within Budget'
  end as budget_status,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_monthly
order by
  usage_date desc,
  pre_tax_cost_amount desc;
```

```sql+sqlite
select
  resource_group,
  usage_date,
  pre_tax_cost_amount,
  case
    when pre_tax_cost_amount > 1000 then 'Over Budget'
    when pre_tax_cost_amount > 800 then 'Near Budget'
    else 'Within Budget'
  end as budget_status,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_monthly
order by
  usage_date desc,
  pre_tax_cost_amount desc;
```

### Resource group cost distribution

Analyze the distribution of costs across resource groups to understand how spending is allocated across different teams or projects.

```sql+postgres
select
  resource_group,
  avg(pre_tax_cost_amount) as avg_monthly_cost,
  min(pre_tax_cost_amount) as min_monthly_cost,
  max(pre_tax_cost_amount) as max_monthly_cost,
  stddev(pre_tax_cost_amount) as cost_variability,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_monthly
group by
  resource_group,
  pre_tax_cost_unit
order by
  avg_monthly_cost desc;
```

```sql+sqlite
select
  resource_group,
  avg(pre_tax_cost_amount) as avg_monthly_cost,
  min(pre_tax_cost_amount) as min_monthly_cost,
  max(pre_tax_cost_amount) as max_monthly_cost,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_monthly
group by
  resource_group,
  pre_tax_cost_unit
order by
  avg_monthly_cost desc;
```

### Quarter-over-quarter cost comparison

Compare quarterly costs for resource groups to understand seasonal patterns and long-term trends.

```sql+postgres
select
  resource_group,
  extract(year from usage_date) as year,
  extract(quarter from usage_date) as quarter,
  sum(pre_tax_cost_amount) as quarterly_cost,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_monthly
group by
  resource_group,
  extract(year from usage_date),
  extract(quarter from usage_date),
  pre_tax_cost_unit
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
  sum(pre_tax_cost_amount) as quarterly_cost,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_monthly
group by
  resource_group,
  cast(strftime('%Y', usage_date) as integer),
  case
    when cast(strftime('%m', usage_date) as integer) <= 3 then 1
    when cast(strftime('%m', usage_date) as integer) <= 6 then 2
    when cast(strftime('%m', usage_date) as integer) <= 9 then 3
    else 4
  end,
  pre_tax_cost_unit
order by
  resource_group,
  year,
  quarter;
```

### Compare pre-tax vs amortized costs

Analyze the difference between pre-tax costs and amortized costs to understand reservation impacts.

```sql+postgres
select
  resource_group,
  usage_date,
  pre_tax_cost_amount,
  amortized_cost_amount,
  (pre_tax_cost_amount - amortized_cost_amount) as reservation_savings,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_monthly
where
  amortized_cost_amount is not null
  and pre_tax_cost_amount != amortized_cost_amount
order by
  reservation_savings desc;
```

```sql+sqlite
select
  resource_group,
  usage_date,
  pre_tax_cost_amount,
  amortized_cost_amount,
  (pre_tax_cost_amount - amortized_cost_amount) as reservation_savings,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_monthly
where
  amortized_cost_amount is not null
  and pre_tax_cost_amount != amortized_cost_amount
order by
  reservation_savings desc;
```

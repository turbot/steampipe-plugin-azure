---
title: "Steampipe Table: azure_cost_by_resource_group_monthly - Query Azure Monthly Resource Group Costs using SQL"
description: "Allows users to query Azure Monthly Resource Group Costs, providing detailed cost breakdown by resource group on a monthly basis."
folder: "Cost Management"
---

# Table: azure_cost_by_resource_group_monthly - Query Azure Monthly Resource Group Costs using SQL

Azure Cost Management provides cost analytics to help you understand and manage your Azure spending. The monthly resource group cost breakdown allows you to track costs at a higher level, providing insights into how much each Azure resource group costs on a month-by-month basis. This helps in understanding long-term cost trends per team or project, budget planning, and strategic cost management decisions for different departments.

## Table Usage Guide

The `azure_cost_by_resource_group_monthly` table provides insights into monthly cost breakdown by resource group within Microsoft Azure. As a Cloud Architect, FinOps engineer, or DevOps professional, explore resource group-specific monthly cost details through this table, including monthly usage costs, currency information, and resource group names. Utilize it to uncover monthly cost patterns per team or project, compare resource group costs across months, track long-term spending trends, and support departmental budget planning and forecasting.

## Examples

### Basic monthly cost info by resource group

Explore monthly costs across different Azure resource groups to understand your spending patterns and identify the most expensive resource groups by month.

```sql+postgres
select
  usage_date,
  resource_group,
  unblended_cost_amount,
  unblended_cost_unit
from
  azure_cost_by_resource_group_monthly
order by
  usage_date desc,
  unblended_cost_amount desc;
```

```sql+sqlite
select
  usage_date,
  resource_group,
  unblended_cost_amount,
  unblended_cost_unit
from
  azure_cost_by_resource_group_monthly
order by
  usage_date desc,
  unblended_cost_amount desc;
```

### Monthly cost trend for a specific resource group

Analyze the monthly cost trend for a specific Azure resource group to understand its long-term usage patterns and cost evolution.

```sql+postgres
select
  usage_date,
  resource_group,
  unblended_cost_amount,
  unblended_cost_unit
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
  unblended_cost_amount,
  unblended_cost_unit
from
  azure_cost_by_resource_group_monthly
where
  resource_group = 'production-rg'
order by
  usage_date desc;
```

### Total monthly spend by resource group

Get the total monthly cost for each resource group to understand which teams or projects contribute most to your monthly Azure bill.

```sql+postgres
select
  resource_group,
  sum(unblended_cost_amount) as total_monthly_cost,
  unblended_cost_unit,
  count(*) as months_with_usage
from
  azure_cost_by_resource_group_monthly
group by
  resource_group,
  unblended_cost_unit
order by
  total_monthly_cost desc;
```

```sql+sqlite
select
  resource_group,
  sum(unblended_cost_amount) as total_monthly_cost,
  unblended_cost_unit,
  count(*) as months_with_usage
from
  azure_cost_by_resource_group_monthly
group by
  resource_group,
  unblended_cost_unit
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
    unblended_cost_amount,
    lag(unblended_cost_amount) over (partition by resource_group order by usage_date) as prev_month_cost
  from
    azure_cost_by_resource_group_monthly
)
select
  resource_group,
  usage_date,
  unblended_cost_amount as current_month,
  prev_month_cost as previous_month,
  unblended_cost_amount - prev_month_cost as cost_change
from
  monthly_costs
where
  prev_month_cost is not null
  and unblended_cost_amount > prev_month_cost
order by
  cost_change desc;
```

```sql+sqlite
with monthly_costs as (
  select
    resource_group,
    usage_date,
    unblended_cost_amount,
    lag(unblended_cost_amount) over (partition by resource_group order by usage_date) as prev_month_cost
  from
    azure_cost_by_resource_group_monthly
)
select
  resource_group,
  usage_date,
  unblended_cost_amount as current_month,
  prev_month_cost as previous_month,
  unblended_cost_amount - prev_month_cost as cost_change
from
  monthly_costs
where
  prev_month_cost is not null
  and unblended_cost_amount > prev_month_cost
order by
  cost_change desc;
```

### Departmental budget tracking

Track monthly spending against potential budgets by resource group to monitor departmental or team spending limits.

```sql+postgres
select
  resource_group,
  usage_date,
  unblended_cost_amount,
  case
    when unblended_cost_amount > 1000 then 'Over Budget'
    when unblended_cost_amount > 800 then 'Near Budget'
    else 'Within Budget'
  end as budget_status,
  unblended_cost_unit
from
  azure_cost_by_resource_group_monthly
order by
  usage_date desc,
  unblended_cost_amount desc;
```

```sql+sqlite
select
  resource_group,
  usage_date,
  unblended_cost_amount,
  case
    when unblended_cost_amount > 1000 then 'Over Budget'
    when unblended_cost_amount > 800 then 'Near Budget'
    else 'Within Budget'
  end as budget_status,
  unblended_cost_unit
from
  azure_cost_by_resource_group_monthly
order by
  usage_date desc,
  unblended_cost_amount desc;
```

### Resource group cost distribution

Analyze the distribution of costs across resource groups to understand how spending is allocated across different teams or projects.

```sql+postgres
select
  resource_group,
  avg(unblended_cost_amount) as avg_monthly_cost,
  min(unblended_cost_amount) as min_monthly_cost,
  max(unblended_cost_amount) as max_monthly_cost,
  stddev(unblended_cost_amount) as cost_variability,
  unblended_cost_unit
from
  azure_cost_by_resource_group_monthly
group by
  resource_group,
  unblended_cost_unit
order by
  avg_monthly_cost desc;
```

```sql+sqlite
select
  resource_group,
  avg(unblended_cost_amount) as avg_monthly_cost,
  min(unblended_cost_amount) as min_monthly_cost,
  max(unblended_cost_amount) as max_monthly_cost,
  unblended_cost_unit
from
  azure_cost_by_resource_group_monthly
group by
  resource_group,
  unblended_cost_unit
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
  sum(unblended_cost_amount) as quarterly_cost,
  unblended_cost_unit
from
  azure_cost_by_resource_group_monthly
group by
  resource_group,
  extract(year from usage_date),
  extract(quarter from usage_date),
  unblended_cost_unit
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
  sum(unblended_cost_amount) as quarterly_cost,
  unblended_cost_unit
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
  unblended_cost_unit
order by
  resource_group,
  year,
  quarter;
```

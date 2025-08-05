---
title: "Steampipe Table: azure_cost_by_resource_group_daily - Query Azure Daily Resource Group Costs using SQL"
description: "Allows users to query Azure Daily Resource Group Costs, providing detailed cost breakdown by resource group on a daily basis."
folder: "Cost Management"
---

# Table: azure_cost_by_resource_group_daily - Query Azure Daily Resource Group Costs using SQL

Azure Cost Management provides cost analytics to help you understand and manage your Azure spending. The daily resource group cost breakdown allows you to track costs at a granular level, providing insights into how much each Azure resource group costs on a day-by-day basis. This helps in identifying cost trends per team or project, optimizing resource allocation, and managing departmental budgets effectively.

## Table Usage Guide

The `azure_cost_by_resource_group_daily` table provides insights into daily cost breakdown by resource group within Microsoft Azure. As a Cloud Architect, FinOps engineer, or DevOps professional, explore resource group-specific cost details through this table, including daily usage costs, currency information, and resource group names. Utilize it to uncover cost patterns per team or project, identify expensive resource groups, track daily spending trends, and optimize resource group allocation.

**Important Notes:**

- This table supports optional quals. Queries with optional quals are optimised to reduce query time and improve performance. Optional quals are supported for the following columns:
  - `scope` with supported operators `=`.
  - `type` with supported operators `=`. Valid values are 'ActualCost' (default) and 'AmortizedCost'.
  - `period_start` with supported operators `=`, `>=`, `>`, `<=`, and `<`.
  - `period_end` with supported operators `=`, `>=`, `>`, `<=`, and `<`.

## Examples

### Basic daily cost info by resource group

Explore daily costs across different Azure resource groups to understand your spending patterns and identify the most expensive resource groups.

```sql+postgres
select
  usage_date,
  resource_group,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_daily
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
  azure_cost_by_resource_group_daily
order by
  usage_date desc,
  pre_tax_cost_amount desc;
```

### Daily costs for a specific resource group

Analyze the daily cost trend for a specific Azure resource group to understand its usage patterns and cost fluctuations.

```sql+postgres
select
  usage_date,
  resource_group,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_daily
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
  azure_cost_by_resource_group_daily
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
  azure_cost_by_resource_group_daily
where
  period_start = '2024-08-01'
  and period_end = '2024-08-31'
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
  azure_cost_by_resource_group_daily
where
  period_start = '2024-08-01'
  and period_end = '2024-08-31'
order by
  pre_tax_cost_amount desc;
```

### Top 5 most expensive resource groups yesterday

Identify the most expensive Azure resource groups from the previous day to focus cost optimization efforts on specific teams or projects.

```sql+postgres
select
  resource_group,
  sum(pre_tax_cost_amount) as total_cost,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_daily
where
  usage_date >= current_date - interval '1 day'
  and usage_date < current_date
group by
  resource_group,
  pre_tax_cost_unit
order by
  total_cost desc
limit 5;
```

```sql+sqlite
select
  resource_group,
  sum(pre_tax_cost_amount) as total_cost,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_daily
where
  usage_date >= date('now', '-1 day')
  and usage_date < date('now')
group by
  resource_group,
  pre_tax_cost_unit
order by
  total_cost desc
limit 5;
```

### Weekly cost trend for all resource groups

Analyze the weekly cost trend to understand spending patterns across resource groups and identify cost spikes.

```sql+postgres
select
  date_trunc('week', usage_date) as week_start,
  sum(pre_tax_cost_amount) as weekly_cost,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_daily
where
  usage_date >= current_date - interval '30 days'
group by
  date_trunc('week', usage_date),
  pre_tax_cost_unit
order by
  week_start desc;
```

```sql+sqlite
select
  date(usage_date, 'weekday 0', '-6 days') as week_start,
  sum(pre_tax_cost_amount) as weekly_cost,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_daily
where
  usage_date >= date('now', '-30 days')
group by
  date(usage_date, 'weekday 0', '-6 days'),
  pre_tax_cost_unit
order by
  week_start desc;
```

### Resource groups with zero costs

Find resource groups that had no costs on specific days, which might indicate unused or idle resources that could be cleaned up.

```sql+postgres
select
  usage_date,
  resource_group,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_daily
where
  pre_tax_cost_amount = 0
order by
  usage_date desc,
  resource_group;
```

```sql+sqlite
select
  usage_date,
  resource_group,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_daily
where
  pre_tax_cost_amount = 0
order by
  usage_date desc,
  resource_group;
```

### Daily cost comparison between resource groups

Compare daily costs between different resource groups to understand relative spending across teams or projects.

```sql+postgres
select
  usage_date,
  resource_group,
  pre_tax_cost_amount,
  pre_tax_cost_unit,
  rank() over (partition by usage_date order by pre_tax_cost_amount desc) as cost_rank
from
  azure_cost_by_resource_group_daily
where
  usage_date >= current_date - interval '7 days'
order by
  usage_date desc,
  cost_rank;
```

```sql+sqlite
select
  usage_date,
  resource_group,
  pre_tax_cost_amount,
  pre_tax_cost_unit,
  rank() over (partition by usage_date order by pre_tax_cost_amount desc) as cost_rank
from
  azure_cost_by_resource_group_daily
where
  usage_date >= date('now', '-7 days')
order by
  usage_date desc,
  cost_rank;
```

### Resource groups with costs above threshold

Find resource groups that exceeded a specific cost threshold on any given day for cost monitoring and budget alerting.

```sql+postgres
select
  usage_date,
  resource_group,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_daily
where
  pre_tax_cost_amount > 50
order by
  pre_tax_cost_amount desc;
```

```sql+sqlite
select
  usage_date,
  resource_group,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_resource_group_daily
where
  pre_tax_cost_amount > 50
order by
  pre_tax_cost_amount desc;
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
  azure_cost_by_resource_group_daily
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
  azure_cost_by_resource_group_daily
where
  amortized_cost_amount is not null
  and pre_tax_cost_amount != amortized_cost_amount
order by
  reservation_savings desc;
```

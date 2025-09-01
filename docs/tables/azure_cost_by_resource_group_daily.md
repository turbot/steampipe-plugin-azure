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

- You **_must_** specify `cost_type` (ActualCost or AmortizedCost) in a `where` clause in order to use this table.
- For improved performance, it is advised that you use the optional quals `period_start` and `period_end` to limit the result set to a specific time period.
- This table supports optional quals. Queries with optional quals are optimised to use Azure Cost Management filters. Optional quals are supported for the following columns:
  - `scope` with supported operators `=`. Default to current subscription. Possible value are see: [Supported Scope](https://learn.microsoft.com/en-gb/rest/api/cost-management/query/usage?view=rest-cost-management-2025-03-01&tabs=HTTP#uri-parameters)
  - `period_start` with supported operators `=`. Default: 1 year ago.
  - `period_end` with supported operators `=`. Default: yesterday.
  - `resource_group` with supported operators `=`, `<>`.

## Examples

### Recent daily costs by resource group
Get the most recent 7 days of daily costs across Azure resource groups, showing the cost breakdown by resource group with subscription details.

```sql+postgres
select
  usage_date,
  resource_group,
  cost,
  currency,
  subscription_id
from
  azure_cost_by_resource_group_daily
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
  usage_date,
  resource_group,
  cost,
  currency,
  subscription_id
from
  azure_cost_by_resource_group_daily
where
  cost_type = 'ActualCost'
  and usage_date >= date('now', '-7 days')
order by
  usage_date desc,
  cost desc
limit 10;
```

### Historical daily costs for a specific resource group
Analyze the complete historical daily cost trend for a specific Azure resource group to understand its usage patterns and cost evolution over time.

```sql+postgres
select
  usage_date,
  resource_group,
  cost,
  pre_tax_cost,
  currency
from
  azure_cost_by_resource_group_daily
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
  azure_cost_by_resource_group_daily
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
  azure_cost_by_resource_group_daily
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
  resource_group,
  cost,
  pre_tax_cost,
  currency,
  period_start,
  period_end
from
  azure_cost_by_resource_group_daily
where
  cost_type = 'ActualCost'
  and period_start = '2025-08-01'
  and period_end = '2025-08-31'
order by
  cost desc;
```

### Top 5 most expensive resource groups in the last 30 days
Identify the most expensive Azure resource groups from the last 30 days to focus cost optimization efforts on specific teams or projects.

```sql+postgres
select
  resource_group,
  sum(cost) as total_cost,
  avg(cost) as avg_daily_cost,
  currency
from
  azure_cost_by_resource_group_daily
where
  cost_type = 'ActualCost'
  and usage_date >= NOW() - INTERVAL '30 days'
group by
  resource_group,
  currency
order by
  total_cost desc
limit 5;
```

```sql+sqlite
select
  resource_group,
  sum(cost) as total_cost,
  avg(cost) as avg_daily_cost,
  currency
from
  azure_cost_by_resource_group_daily
where
  cost_type = 'ActualCost'
  and usage_date >= date('now', '-30 days')
group by
  resource_group,
  currency
order by
  total_cost desc
limit 5;
```

### Daily total spending trends
Analyze the aggregated daily cost trends across all resource groups to identify spending patterns and cost spikes over the last 30 days.

```sql+postgres
select
  usage_date,
  sum(cost) as daily_total_cost,
  currency
from
  azure_cost_by_resource_group_daily
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
  azure_cost_by_resource_group_daily
where
  cost_type = 'ActualCost'
  and usage_date >= date('now', '-30 days')
group by
  usage_date,
  currency
order by
  usage_date desc;
```

### Identify idle resource groups
Find resource groups that had zero costs on specific days, which might indicate unused or idle resources that could be optimized or cleaned up.

```sql+postgres
select
  usage_date,
  resource_group,
  cost,
  currency
from
  azure_cost_by_resource_group_daily
where
  cost_type = 'ActualCost'
  and cost = 0
order by
  usage_date desc,
  resource_group;
```

```sql+sqlite
select
  usage_date,
  resource_group,
  cost,
  currency
from
  azure_cost_by_resource_group_daily
where
  cost_type = 'ActualCost'
  and cost = 0
order by
  usage_date desc,
  resource_group;
```

### Daily cost ranking by resource group
Compare daily costs between different resource groups using ranking to understand relative spending patterns across teams or projects over the last 7 days.

```sql+postgres
select
  usage_date,
  resource_group,
  cost,
  currency,
  rank() over (partition by usage_date order by cost desc) as cost_rank
from
  azure_cost_by_resource_group_daily
where
  cost_type = 'ActualCost'
  and usage_date >= NOW() - INTERVAL '7 days'
order by
  usage_date desc,
  cost_rank;
```

```sql+sqlite
select
  usage_date,
  resource_group,
  cost,
  currency,
  rank() over (partition by usage_date order by cost desc) as cost_rank
from
  azure_cost_by_resource_group_daily
where
  cost_type = 'ActualCost'
  and usage_date >= date('now', '-7 days')
order by
  usage_date desc,
  cost_rank;
```

### High-cost resource groups above threshold
Identify resource groups that exceeded a specific cost threshold ($1.00) on any given day, useful for cost monitoring and budget alerting.

```sql+postgres
select
  usage_date,
  resource_group,
  cost,
  currency
from
  azure_cost_by_resource_group_daily
where
  cost_type = 'ActualCost'
  and cost > 1.0
order by
  cost desc;
```

```sql+sqlite
select
  usage_date,
  resource_group,
  cost,
  currency
from
  azure_cost_by_resource_group_daily
where
  cost_type = 'ActualCost'
  and cost > 1.0
order by
  cost desc;
```

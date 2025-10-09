---
title: "Steampipe Table: azure_cost_forecast_daily - Query Azure Daily Cost Forecasts using SQL"
description: "Allows users to query Azure Daily Cost Forecasts, providing insights into predicted future costs on a daily basis."
folder: "Cost Management"
---

# Table: azure_cost_forecast_daily - Query Azure Daily Cost Forecasts using SQL

Azure Cost Management provides cost forecasting capabilities to help you predict and plan your future Azure spending. The daily cost forecast provides granular predictions for your Azure costs over the next 90 days, helping with short-term budget planning and cost optimization decisions.

## Table Usage Guide

The `azure_cost_forecast_daily` table provides insights into predicted future costs within Microsoft Azure at a daily granularity. As a FinOps engineer, Cloud Architect, or DevOps professional, explore daily cost forecasts through this table, including predicted costs and currency information. Utilize it to understand near-term cost trends, plan daily budgets, and make proactive cost optimization decisions.

**Important Notes:**
- You **_must_** specify `cost_type` (ActualCost or AmortizedCost) in a `where` clause in order to use this table.
- By default, forecasts are generated for the next 90 days from the current date.
- This table supports optional quals. Queries with optional quals are optimised to use Azure Cost Management filters. Optional quals are supported for the following columns:
  - `scope` with supported operators `=`. Default to current subscription. Possible values are see: [Supported Scope](https://learn.microsoft.com/en-gb/rest/api/cost-management/query/usage?view=rest-cost-management-2025-03-01&tabs=HTTP#uri-parameters)
  - `period_start` with supported operators `=`. Default: current date.
  - `period_end` with supported operators `=`. Default: 90 days from current date.

## Examples

### Basic daily cost forecast
Get the predicted costs for the next 30 days, showing the forecast value.

```sql+postgres
select
  usage_date,
  round(pre_tax_cost::numeric, 2) as forecasted_cost,
  currency
from
  azure_cost_forecast_daily
where
  cost_type = 'ActualCost'
  and period_start = current_date
  and period_end = current_date + interval '30 days'
order by
  usage_date;
```

```sql+sqlite
select
  usage_date,
  round(pre_tax_cost, 2) as forecasted_cost,
  currency
from
  azure_cost_forecast_daily
where
  cost_type = 'ActualCost'
  and period_start = date('now')
  and period_end = date('now', '+30 days')
order by
  usage_date;
```

### Forecast for specific time period
Use period_start and period_end parameters to get forecasts for a specific time range.

```sql+postgres
select
  usage_date,
  round(pre_tax_cost::numeric, 2) as forecasted_cost,
  currency,
  period_start,
  period_end
from
  azure_cost_forecast_daily
where
  cost_type = 'ActualCost'
  and period_start = '2025-01-01'
  and period_end = '2025-01-31'
order by
  usage_date;
```

```sql+sqlite
select
  usage_date,
  round(pre_tax_cost, 2) as forecasted_cost,
  currency,
  period_start,
  period_end
from
  azure_cost_forecast_daily
where
  cost_type = 'ActualCost'
  and period_start = '2025-01-01'
  and period_end = '2025-01-31'
order by
  usage_date;
```

### Compare forecast with historical costs
Compare tomorrow's forecast with today's costs.

```sql+postgres
select
  f.usage_date as forecast_date,
  round(f.pre_tax_cost::numeric, 2) as forecasted_cost,
  round(h.pre_tax_cost::numeric, 2) as current_day_cost,
  round((f.pre_tax_cost - h.pre_tax_cost)::numeric, 2) as cost_difference,
  round(((f.pre_tax_cost - h.pre_tax_cost) / nullif(h.pre_tax_cost, 0) * 100)::numeric, 2) as percentage_change,
  f.currency
from
  azure_cost_forecast_daily f
  left join azure_cost_by_service_daily h on
    date_trunc('day', f.usage_date) = date_trunc('day', current_date + interval '1 day')
where
  f.cost_type = 'ActualCost'
  and h.cost_type = 'ActualCost'
  and f.period_start = current_date
  and f.period_end = current_date + interval '1 day'
order by
  f.usage_date;
```

```sql+sqlite
select
  f.usage_date as forecast_date,
  round(f.pre_tax_cost, 2) as forecasted_cost,
  round(h.pre_tax_cost, 2) as current_day_cost,
  round((f.pre_tax_cost - h.pre_tax_cost), 2) as cost_difference,
  round(((f.pre_tax_cost - h.pre_tax_cost) / nullif(h.pre_tax_cost, 0) * 100), 2) as percentage_change,
  f.currency
from
  azure_cost_forecast_daily f
  left join azure_cost_by_service_daily h on
    date(f.usage_date) = date('now', '+1 day')
where
  f.cost_type = 'ActualCost'
  and h.cost_type = 'ActualCost'
  and f.period_start = date('now')
  and f.period_end = date('now', '+1 day')
order by
  f.usage_date;
```

### Resource group level forecasts for next week
Get cost forecasts broken down by resource group for the next 7 days.

```sql+postgres
select
  scope,
  usage_date,
  round(pre_tax_cost::numeric, 2) as forecasted_cost,
  currency
from
  azure_cost_forecast_daily
where
  cost_type = 'ActualCost'
  and scope like '/subscriptions/%/resourceGroups/%'
  and period_start = current_date
  and period_end = current_date + interval '7 days'
order by
  pre_tax_cost desc;
```

```sql+sqlite
select
  scope,
  usage_date,
  round(pre_tax_cost, 2) as forecasted_cost,
  currency
from
  azure_cost_forecast_daily
where
  cost_type = 'ActualCost'
  and scope like '/subscriptions/%/resourceGroups/%'
  and period_start = date('now')
  and period_end = date('now', '+7 days')
order by
  pre_tax_cost desc;
```

### Weekly forecast aggregation
Aggregate daily forecasts into weekly totals for easier trend analysis.

```sql+postgres
select
  date_trunc('week', usage_date) as week_start,
  round(sum(pre_tax_cost)::numeric, 2) as weekly_forecast,
  round(avg(pre_tax_cost)::numeric, 2) as avg_daily_forecast,
  currency
from
  azure_cost_forecast_daily
where
  cost_type = 'ActualCost'
  and period_start = current_date
  and period_end = current_date + interval '30 days'
group by
  date_trunc('week', usage_date),
  currency
order by
  week_start;
```

```sql+sqlite
select
  strftime('%Y-%m-%d', date(usage_date, 'weekday 0', '-7 days')) as week_start,
  round(sum(pre_tax_cost), 2) as weekly_forecast,
  round(avg(pre_tax_cost), 2) as avg_daily_forecast,
  currency
from
  azure_cost_forecast_daily
where
  cost_type = 'ActualCost'
  and period_start = date('now')
  and period_end = date('now', '+30 days')
group by
  strftime('%Y-%W', usage_date),
  currency
order by
  week_start;
```

### Daily forecast trend analysis for next month
Analyze how the forecast changes over time by comparing consecutive days.

```sql+postgres
select
  usage_date,
  round(pre_tax_cost::numeric, 2) as forecasted_cost,
  round(lag(pre_tax_cost) over (order by usage_date)::numeric, 2) as previous_day_forecast,
  round(((pre_tax_cost - lag(pre_tax_cost) over (order by usage_date)) / 
    nullif(lag(pre_tax_cost) over (order by usage_date), 0) * 100)::numeric, 2) 
    as day_over_day_change
from
  azure_cost_forecast_daily
where
  cost_type = 'ActualCost'
  and period_start = current_date
  and period_end = current_date + interval '30 days'
order by
  usage_date;
```

```sql+sqlite
select
  usage_date,
  round(pre_tax_cost, 2) as forecasted_cost,
  round(lag(pre_tax_cost) over (order by usage_date), 2) as previous_day_forecast,
  round(((pre_tax_cost - lag(pre_tax_cost) over (order by usage_date)) / 
    nullif(lag(pre_tax_cost) over (order by usage_date), 0) * 100), 2) 
    as day_over_day_change
from
  azure_cost_forecast_daily
where
  cost_type = 'ActualCost'
  and period_start = date('now')
  and period_end = date('now', '+30 days')
order by
  usage_date;
```

### Weekend vs weekday cost forecast comparison
Compare forecasted costs between weekdays and weekends to identify usage patterns.

```sql+postgres
select
  case
    when extract(dow from usage_date) in (0, 6) then 'Weekend'
    else 'Weekday'
  end as day_type,
  round(avg(pre_tax_cost)::numeric, 2) as avg_forecast,
  currency
from
  azure_cost_forecast_daily
where
  cost_type = 'ActualCost'
  and period_start = current_date
  and period_end = current_date + interval '30 days'
group by
  day_type,
  currency
order by
  day_type;
```

```sql+sqlite
select
  case
    when cast(strftime('%w', usage_date) as integer) in (0, 6) then 'Weekend'
    else 'Weekday'
  end as day_type,
  round(avg(pre_tax_cost), 2) as avg_forecast,
  currency
from
  azure_cost_forecast_daily
where
  cost_type = 'ActualCost'
  and period_start = date('now')
  and period_end = date('now', '+30 days')
group by
  day_type,
  currency
order by
  day_type;
```
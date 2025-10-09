---
title: "Steampipe Table: azure_cost_forecast_monthly - Query Azure Monthly Cost Forecasts using SQL"
description: "Allows users to query Azure Monthly Cost Forecasts, providing insights into predicted future costs on a monthly basis."
folder: "Cost Management"
---

# Table: azure_cost_forecast_monthly - Query Azure Monthly Cost Forecasts using SQL

Azure Cost Management provides cost forecasting capabilities to help you predict and plan your future Azure spending. The monthly cost forecast provides predictions for your Azure costs over the next year, helping with budget planning and cost optimization decisions.

## Table Usage Guide

The `azure_cost_forecast_monthly` table provides insights into predicted future costs within Microsoft Azure. As a FinOps engineer, Cloud Architect, or DevOps professional, explore monthly cost forecasts through this table, including predicted costs and currency information. Utilize it to understand future cost trends, plan budgets, and make proactive cost optimization decisions.

**Important Notes:**
- You **_must_** specify `cost_type` (ActualCost or AmortizedCost) in a `where` clause in order to use this table.
- By default, forecasts are generated for the next 12 months from the current date.
- This table supports optional quals. Queries with optional quals are optimised to use Azure Cost Management filters. Optional quals are supported for the following columns:
  - `scope` with supported operators `=`. Default to current subscription. Possible values are see: [Supported Scope](https://learn.microsoft.com/en-gb/rest/api/cost-management/query/usage?view=rest-cost-management-2025-03-01&tabs=HTTP#uri-parameters)
  - `period_start` with supported operators `=`. Default: current date.
  - `period_end` with supported operators `=`. Default: 12 months from current date.

## Examples

### Basic monthly cost forecast
Get the predicted costs for the next 6 months, showing the forecast value.

```sql+postgres
select
  usage_date,
  round(pre_tax_cost::numeric, 2) as forecasted_cost,
  currency
from
  azure_cost_forecast_monthly
where
  cost_type = 'ActualCost'
  and period_start = current_date
  and period_end = current_date + interval '6 months'
order by
  usage_date;
```

```sql+sqlite
select
  usage_date,
  round(pre_tax_cost, 2) as forecasted_cost,
  currency
from
  azure_cost_forecast_monthly
where
  cost_type = 'ActualCost'
  and period_start = date('now')
  and period_end = date('now', '+6 months')
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
  azure_cost_forecast_monthly
where
  cost_type = 'ActualCost'
  and period_start = '2025-01-01'
  and period_end = '2025-03-31'
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
  azure_cost_forecast_monthly
where
  cost_type = 'ActualCost'
  and period_start = '2025-01-01'
  and period_end = '2025-03-31'
order by
  usage_date;
```

### Compare forecast with historical costs
Compare next month's forecast with current month's costs.

```sql+postgres
select
  f.usage_date as forecast_month,
  round(f.pre_tax_cost::numeric, 2) as forecasted_cost,
  round(h.pre_tax_cost::numeric, 2) as current_month_cost,
  round((f.pre_tax_cost - h.pre_tax_cost)::numeric, 2) as cost_difference,
  round(((f.pre_tax_cost - h.pre_tax_cost) / nullif(h.pre_tax_cost, 0) * 100)::numeric, 2) as percentage_change,
  f.currency
from
  azure_cost_forecast_monthly f
  left join azure_cost_by_service_monthly h on
    date_trunc('month', f.usage_date) = date_trunc('month', current_date + interval '1 month')
where
  f.cost_type = 'ActualCost'
  and h.cost_type = 'ActualCost'
  and f.period_start = current_date
  and f.period_end = current_date + interval '1 month'
order by
  f.usage_date;
```

```sql+sqlite
select
  f.usage_date as forecast_month,
  round(f.pre_tax_cost, 2) as forecasted_cost,
  round(h.pre_tax_cost, 2) as current_month_cost,
  round((f.pre_tax_cost - h.pre_tax_cost), 2) as cost_difference,
  round(((f.pre_tax_cost - h.pre_tax_cost) / nullif(h.pre_tax_cost, 0) * 100), 2) as percentage_change,
  f.currency
from
  azure_cost_forecast_monthly f
  left join azure_cost_by_service_monthly h on
    strftime('%Y-%m', f.usage_date) = strftime('%Y-%m', date('now', '+1 month'))
where
  f.cost_type = 'ActualCost'
  and h.cost_type = 'ActualCost'
  and f.period_start = date('now')
  and f.period_end = date('now', '+1 month')
order by
  f.usage_date;
```

### Resource group level forecasts for next quarter
Get cost forecasts broken down by resource group for the next 3 months.

```sql+postgres
select
  scope,
  usage_date,
  round(pre_tax_cost::numeric, 2) as forecasted_cost,
  currency
from
  azure_cost_forecast_monthly
where
  cost_type = 'ActualCost'
  and scope like '/subscriptions/%/resourceGroups/%'
  and period_start = current_date
  and period_end = current_date + interval '3 months'
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
  azure_cost_forecast_monthly
where
  cost_type = 'ActualCost'
  and scope like '/subscriptions/%/resourceGroups/%'
  and period_start = date('now')
  and period_end = date('now', '+3 months')
order by
  pre_tax_cost desc;
```

### Monthly forecast trend analysis for next 6 months
Analyze how the forecast changes over time by comparing consecutive months.

```sql+postgres
select
  usage_date,
  round(pre_tax_cost::numeric, 2) as forecasted_cost,
  round(lag(pre_tax_cost) over (order by usage_date)::numeric, 2) as previous_month_forecast,
  round(((pre_tax_cost - lag(pre_tax_cost) over (order by usage_date)) / 
    nullif(lag(pre_tax_cost) over (order by usage_date), 0) * 100)::numeric, 2) 
    as month_over_month_change
from
  azure_cost_forecast_monthly
where
  cost_type = 'ActualCost'
  and period_start = current_date
  and period_end = current_date + interval '6 months'
order by
  usage_date;
```

```sql+sqlite
select
  usage_date,
  round(pre_tax_cost, 2) as forecasted_cost,
  round(lag(pre_tax_cost) over (order by usage_date), 2) as previous_month_forecast,
  round(((pre_tax_cost - lag(pre_tax_cost) over (order by usage_date)) / 
    nullif(lag(pre_tax_cost) over (order by usage_date), 0) * 100), 2) 
    as month_over_month_change
from
  azure_cost_forecast_monthly
where
  cost_type = 'ActualCost'
  and period_start = date('now')
  and period_end = date('now', '+6 months')
order by
  usage_date;
```

### Subscription cost forecast summary for next quarter
Get a summary of forecasted costs across all subscriptions for the next 3 months.

```sql+postgres
select
  subscription_id,
  min(usage_date) as forecast_start,
  max(usage_date) as forecast_end,
  round(avg(pre_tax_cost)::numeric, 2) as avg_monthly_forecast,
  round(sum(pre_tax_cost)::numeric, 2) as total_forecast,
  currency
from
  azure_cost_forecast_monthly
where
  cost_type = 'ActualCost'
  and period_start = current_date
  and period_end = current_date + interval '3 months'
group by
  subscription_id,
  currency
order by
  total_forecast desc;
```

```sql+sqlite
select
  subscription_id,
  min(usage_date) as forecast_start,
  max(usage_date) as forecast_end,
  round(avg(pre_tax_cost), 2) as avg_monthly_forecast,
  round(sum(pre_tax_cost), 2) as total_forecast,
  currency
from
  azure_cost_forecast_monthly
where
  cost_type = 'ActualCost'
  and period_start = date('now')
  and period_end = date('now', '+3 months')
group by
  subscription_id,
  currency
order by
  total_forecast desc;
```
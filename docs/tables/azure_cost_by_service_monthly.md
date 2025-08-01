---
title: "Steampipe Table: azure_cost_by_service_monthly - Query Azure Monthly Service Costs using SQL"
description: "Allows users to query Azure Monthly Service Costs, providing detailed cost breakdown by service on a monthly basis."
folder: "Cost Management"
---

# Table: azure_cost_by_service_monthly - Query Azure Monthly Service Costs using SQL

Azure Cost Management provides cost analytics to help you understand and manage your Azure spending. The monthly service cost breakdown allows you to track costs at a higher level, providing insights into how much each Azure service costs on a month-by-month basis. This helps in understanding long-term cost trends, budget planning, and strategic cost management decisions.

## Table Usage Guide

The `azure_cost_by_service_monthly` table provides insights into monthly cost breakdown by service within Microsoft Azure. As a Cloud Architect, FinOps engineer, or DevOps professional, explore service-specific monthly cost details through this table, including monthly usage costs, currency information, and service names. Utilize it to uncover monthly cost patterns, compare service costs across months, track long-term spending trends, and support budget planning and forecasting.

## Examples

### Basic monthly cost info

Explore monthly costs across different Azure services to understand your spending patterns and identify the most expensive services by month.

```sql+postgres
select
  usage_date,
  service_name,
  unblended_cost_amount,
  unblended_cost_unit
from
  azure_cost_by_service_monthly
order by
  usage_date desc,
  unblended_cost_amount desc;
```

```sql+sqlite
select
  usage_date,
  service_name,
  unblended_cost_amount,
  unblended_cost_unit
from
  azure_cost_by_service_monthly
order by
  usage_date desc,
  unblended_cost_amount desc;
```

### Monthly cost trend for a specific service

Analyze the monthly cost trend for a specific Azure service to understand its long-term usage patterns and cost evolution.

```sql+postgres
select
  usage_date,
  service_name,
  unblended_cost_amount,
  unblended_cost_unit
from
  azure_cost_by_service_monthly
where
  service_name = 'Storage'
order by
  usage_date desc;
```

```sql+sqlite
select
  usage_date,
  service_name,
  unblended_cost_amount,
  unblended_cost_unit
from
  azure_cost_by_service_monthly
where
  service_name = 'Storage'
order by
  usage_date desc;
```

### Total monthly spend by service

Get the total monthly cost for each service to understand which services contribute most to your monthly Azure bill.

```sql+postgres
select
  service_name,
  sum(unblended_cost_amount) as total_monthly_cost,
  unblended_cost_unit,
  count(*) as months_with_usage
from
  azure_cost_by_service_monthly
group by
  service_name,
  unblended_cost_unit
order by
  total_monthly_cost desc;
```

```sql+sqlite
select
  service_name,
  sum(unblended_cost_amount) as total_monthly_cost,
  unblended_cost_unit,
  count(*) as months_with_usage
from
  azure_cost_by_service_monthly
group by
  service_name,
  unblended_cost_unit
order by
  total_monthly_cost desc;
```

### Year-over-year cost comparison

Compare monthly costs between different years to understand cost growth and seasonal patterns.

```sql+postgres
select
  extract(month from usage_date) as month_number,
  extract(year from usage_date) as year,
  sum(unblended_cost_amount) as monthly_total,
  unblended_cost_unit
from
  azure_cost_by_service_monthly
where
  extract(year from usage_date) in (2024, 2025)
group by
  extract(month from usage_date),
  extract(year from usage_date),
  unblended_cost_unit
order by
  month_number,
  year;
```

```sql+sqlite
select
  cast(strftime('%m', usage_date) as integer) as month_number,
  cast(strftime('%Y', usage_date) as integer) as year,
  sum(unblended_cost_amount) as monthly_total,
  unblended_cost_unit
from
  azure_cost_by_service_monthly
where
  cast(strftime('%Y', usage_date) as integer) in (2024, 2025)
group by
  cast(strftime('%m', usage_date) as integer),
  cast(strftime('%Y', usage_date) as integer),
  unblended_cost_unit
order by
  month_number,
  year;
```

### Services with increasing monthly costs

Identify services where costs are trending upward by comparing recent months to help focus cost optimization efforts.

```sql+postgres
with monthly_costs as (
  select
    service_name,
    usage_date,
    unblended_cost_amount,
    lag(unblended_cost_amount) over (partition by service_name order by usage_date) as prev_month_cost
  from
    azure_cost_by_service_monthly
)
select
  service_name,
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
    service_name,
    usage_date,
    unblended_cost_amount,
    lag(unblended_cost_amount) over (partition by service_name order by usage_date) as prev_month_cost
  from
    azure_cost_by_service_monthly
)
select
  service_name,
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

### Average monthly cost per service

Calculate the average monthly cost for each service to understand baseline spending patterns.

```sql+postgres
select
  service_name,
  round(avg(unblended_cost_amount), 2) as avg_monthly_cost,
  round(min(unblended_cost_amount), 2) as min_monthly_cost,
  round(max(unblended_cost_amount), 2) as max_monthly_cost,
  unblended_cost_unit
from
  azure_cost_by_service_monthly
group by
  service_name,
  unblended_cost_unit
order by
  avg_monthly_cost desc;
```

```sql+sqlite
select
  service_name,
  round(avg(unblended_cost_amount), 2) as avg_monthly_cost,
  round(min(unblended_cost_amount), 2) as min_monthly_cost,
  round(max(unblended_cost_amount), 2) as max_monthly_cost,
  unblended_cost_unit
from
  azure_cost_by_service_monthly
group by
  service_name,
  unblended_cost_unit
order by
  avg_monthly_cost desc;
```

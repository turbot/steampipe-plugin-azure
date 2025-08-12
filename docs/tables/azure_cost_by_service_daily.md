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

- This table supports optional quals. Queries with optional quals are optimised to reduce query time and improve performance. Optional quals are supported for the following columns:
  - `scope` with supported operators `=`.
  - `type` with supported operators `=`. Valid values are 'ActualCost' (default) and 'AmortizedCost'.
  - `period_start` with supported operators `=`, `>=`, `>`, `<=`, and `<`.
  - `period_end` with supported operators `=`, `>=`, `>`, `<=`, and `<`.

## Examples

### Basic daily cost info

Explore daily costs across different Azure services to understand your spending patterns and identify the most expensive services.

```sql+postgres
select
  usage_date,
  service_name,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_service_daily
order by
  usage_date desc,
  pre_tax_cost_amount desc;
```

```sql+sqlite
select
  usage_date,
  service_name,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_service_daily
order by
  usage_date desc,
  pre_tax_cost_amount desc;
```

### Daily costs for a specific service

Analyze the daily cost trend for a specific Azure service to understand its usage patterns and cost fluctuations.

```sql+postgres
select
  usage_date,
  service_name,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_service_daily
where
  service_name = 'Storage'
order by
  usage_date desc;
```

```sql+sqlite
select
  usage_date,
  service_name,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_service_daily
where
  service_name = 'Storage'
order by
  usage_date desc;
```

### Query costs for a specific period

Use period_start and period_end parameters to query costs for a specific time range.

```sql+postgres
select
  usage_date,
  service_name,
  pre_tax_cost_amount,
  pre_tax_cost_unit,
  period_start,
  period_end
from
  azure_cost_by_service_daily
where
  period_start = '2024-08-01'
  and period_end = '2024-08-31'
order by
  pre_tax_cost_amount desc;
```

```sql+sqlite
select
  usage_date,
  service_name,
  pre_tax_cost_amount,
  pre_tax_cost_unit,
  period_start,
  period_end
from
  azure_cost_by_service_daily
where
  period_start = '2024-08-01'
  and period_end = '2024-08-31'
order by
  pre_tax_cost_amount desc;
```

### Top 5 most expensive services yesterday

Identify the most expensive Azure services from the previous day to focus cost optimization efforts.

```sql+postgres
select
  service_name,
  sum(pre_tax_cost_amount) as total_cost,
  pre_tax_cost_unit
from
  azure_cost_by_service_daily
where
  usage_date >= current_date - interval '1 day'
  and usage_date < current_date
group by
  service_name,
  pre_tax_cost_unit
order by
  total_cost desc
limit 5;
```

```sql+sqlite
select
  service_name,
  sum(pre_tax_cost_amount) as total_cost,
  pre_tax_cost_unit
from
  azure_cost_by_service_daily
where
  usage_date >= date('now', '-1 day')
  and usage_date < date('now')
group by
  service_name,
  pre_tax_cost_unit
order by
  total_cost desc
limit 5;
```

### Weekly cost trend for all services

Analyze the weekly cost trend to understand spending patterns and identify cost spikes.

```sql+postgres
select
  date_trunc('week', usage_date) as week_start,
  sum(pre_tax_cost_amount) as weekly_cost,
  pre_tax_cost_unit
from
  azure_cost_by_service_daily
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
  azure_cost_by_service_daily
where
  usage_date >= date('now', '-30 days')
group by
  date(usage_date, 'weekday 0', '-6 days'),
  pre_tax_cost_unit
order by
  week_start desc;
```

### Daily cost by service for a specific date range

Get detailed daily cost breakdown for all services within a specific date range.

```sql+postgres
select
  usage_date,
  service_name,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_service_daily
where
  usage_date between '2025-01-01' and '2025-01-07'
order by
  usage_date,
  pre_tax_cost_amount desc;
```

```sql+sqlite
select
  usage_date,
  service_name,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_service_daily
where
  usage_date between '2025-01-01' and '2025-01-07'
order by
  usage_date,
  pre_tax_cost_amount desc;
```

### Services with cost above threshold

Find services that exceeded a specific cost threshold on any given day for cost monitoring and alerting.

```sql+postgres
select
  usage_date,
  service_name,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_service_daily
where
  pre_tax_cost_amount > 100
order by
  pre_tax_cost_amount desc;
```

```sql+sqlite
select
  usage_date,
  service_name,
  pre_tax_cost_amount,
  pre_tax_cost_unit
from
  azure_cost_by_service_daily
where
  pre_tax_cost_amount > 100
order by
  pre_tax_cost_amount desc;
```

### Compare pre-tax vs amortized costs

Analyze the difference between pre-tax costs and amortized costs to understand reservation impacts.

```sql+postgres
select
  service_name,
  usage_date,
  pre_tax_cost_amount,
  amortized_cost_amount,
  (pre_tax_cost_amount - amortized_cost_amount) as reservation_savings,
  pre_tax_cost_unit
from
  azure_cost_by_service_daily
where
  amortized_cost_amount is not null
  and pre_tax_cost_amount != amortized_cost_amount
order by
  reservation_savings desc;
```

```sql+sqlite
select
  service_name,
  usage_date,
  pre_tax_cost_amount,
  amortized_cost_amount,
  (pre_tax_cost_amount - amortized_cost_amount) as reservation_savings,
  pre_tax_cost_unit
from
  azure_cost_by_service_daily
where
  amortized_cost_amount is not null
  and pre_tax_cost_amount != amortized_cost_amount
order by
  reservation_savings desc;
```

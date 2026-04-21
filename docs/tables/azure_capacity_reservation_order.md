---
title: "Steampipe Table: azure_capacity_reservation_order - Query Azure Capacity Reservation Orders using SQL"
description: "Allows users to query Azure Capacity Reservation Orders, providing details on reservation order lifecycle, billing plans, term durations, and associated reservations."
folder: "Cost Management"
---

# Table: azure_capacity_reservation_order - Query Azure Capacity Reservation Orders using SQL

Azure Capacity Reservation Orders are the top-level container for one or more Azure Capacity Reservations. A reservation order tracks the overall purchase, billing plan (Upfront or Monthly), term (1 or 3  years), and expiration of the reserved capacity commitment. Reservation orders are tenant-level resources and are not associated with any resource group or region.

## Table Usage Guide

The `azure_capacity_reservation_order` table provides insights into reservation orders within Azure. As a FinOps engineer or cloud administrator, use this table to audit reservation orders, track billing plans and expiration timelines, monitor provisioning state, and review the list of individual reservations contained within each order.

## Examples

### Basic info
Get an overview of your Azure capacity reservation orders, including their current provisioning state, billing plan, term length, and the total quantity of SKUs originally purchased.

```sql+postgres
select
  name,
  id,
  display_name,
  provisioning_state,
  billing_plan,
  term,
  original_quantity,
  created_date_time
from
  azure_capacity_reservation_order;
```

```sql+sqlite
select
  name,
  id,
  display_name,
  provisioning_state,
  billing_plan,
  term,
  original_quantity,
  created_date_time
from
  azure_capacity_reservation_order;
```

### List reservation orders expiring within the next 90 days
Identify reservation orders that will expire soon so you can plan renewals or new purchases to maintain continuous reserved capacity coverage.

```sql+postgres
select
  display_name,
  name,
  term,
  billing_plan,
  original_quantity,
  benefit_start_time,
  expiry_date
from
  azure_capacity_reservation_order
where
  expiry_date <= now() + interval '90 days'
order by
  expiry_date;
```

```sql+sqlite
select
  display_name,
  name,
  term,
  billing_plan,
  original_quantity,
  benefit_start_time,
  expiry_date
from
  azure_capacity_reservation_order
where
  expiry_date <= date('now', '+90 days')
order by
  expiry_date;
```

### List reservation orders with monthly billing plans
Find all reservation orders that use a monthly billing plan rather than upfront payment, useful for auditing cash-flow commitments and forecasting monthly reservation charges.

```sql+postgres
select
  display_name,
  name,
  term,
  original_quantity,
  created_date_time,
  expiry_date
from
  azure_capacity_reservation_order
where
  billing_plan = 'Monthly';
```

```sql+sqlite
select
  display_name,
  name,
  term,
  original_quantity,
  created_date_time,
  expiry_date
from
  azure_capacity_reservation_order
where
  billing_plan = 'Monthly';
```

### Count reservations within each order
Understand the size and composition of each reservation order by counting how many individual reservations it contains, which helps identify orders that may need consolidation or further splitting.

```sql+postgres
select
  display_name,
  name,
  term,
  billing_plan,
  jsonb_array_length(reservations) as reservation_count
from
  azure_capacity_reservation_order
order by
  reservation_count desc;
```

```sql+sqlite
select
  display_name,
  name,
  term,
  billing_plan,
  json_array_length(reservations) as reservation_count
from
  azure_capacity_reservation_order
order by
  reservation_count desc;
```

### Get payment schedule details for monthly billing orders
Retrieve the full billing schedule for reservation orders on a monthly plan, enabling detailed review of payment installments and amounts committed over the order's term.

```sql+postgres
select
  display_name,
  name,
  term,
  original_quantity,
  jsonb_pretty(plan_information) as payment_schedule
from
  azure_capacity_reservation_order
where
  billing_plan = 'Monthly';
```

```sql+sqlite
select
  display_name,
  name,
  term,
  original_quantity,
  plan_information as payment_schedule
from
  azure_capacity_reservation_order
where
  billing_plan = 'Monthly';
```

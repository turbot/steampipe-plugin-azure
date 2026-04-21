---
title: "Steampipe Table: azure_capacity_reservation - Query Azure Capacity Reservations using SQL"
description: "Allows users to query Azure Capacity Reservations, providing details on reserved compute capacity, SKU details, billing plans, term durations, and utilization."
folder: "Cost Management"
---

# Table: azure_capacity_reservation - Query Azure Capacity Reservations using SQL

Azure Capacity Reservations allow you to reserve specific VM compute capacity within an Azure region for any duration, ensuring capacity is available when you need it. Each reservation belongs to a reservation order and represents a committed quantity of a specific VM SKU, with billing applied either upfront or monthly over a 1- or 3-year term.

## Table Usage Guide

The `azure_capacity_reservation` table provides insights into individual capacity reservations within Azure. As a cloud architect or FinOps engineer, use this table to audit reserved VM capacity, monitor provisioning states, review applied billing scopes, track expiration dates, and assess utilization of purchased reservations.

## Examples

### Basic info
Explore the essential properties of your Azure capacity reservations to get an overview of what compute capacity has been reserved, its current provisioning state, and associated SKU details.

```sql+postgres
select
  name,
  reservation_order_id,
  reservation_id,
  display_name,
  provisioning_state,
  sku_name,
  quantity,
  region
from
  azure_capacity_reservation;
```

```sql+sqlite
select
  name,
  reservation_order_id,
  reservation_id,
  display_name,
  provisioning_state,
  sku_name,
  quantity,
  region
from
  azure_capacity_reservation;
```

### List reservations expiring within the next 30 days
Identify capacity reservations that are about to expire so you can take action — such as enabling auto-renew or purchasing a new reservation — before capacity commitment ends.

```sql+postgres
select
  display_name,
  reservation_order_id,
  sku_name,
  quantity,
  term,
  expiry_date,
  renew
from
  azure_capacity_reservation
where
  expiry_date <= now() + interval '30 days';
```

```sql+sqlite
select
  display_name,
  reservation_order_id,
  sku_name,
  quantity,
  term,
  expiry_date,
  renew
from
  azure_capacity_reservation
where
  expiry_date <= date('now', '+30 days');
```

### List reservations by applied scope type
Review how reservation discounts are applied across your Azure environment by examining the scope type (Single, Shared, or ManagementGroup) and the specific subscriptions or management groups where the benefit is applied.

```sql+postgres
select
  display_name,
  applied_scope_type,
  applied_scopes,
  billing_scope_id,
  sku_name,
  quantity
from
  azure_capacity_reservation
order by
  applied_scope_type;
```

```sql+sqlite
select
  display_name,
  applied_scope_type,
  applied_scopes,
  billing_scope_id,
  sku_name,
  quantity
from
  azure_capacity_reservation
order by
  applied_scope_type;
```

### List reservations with auto-renew enabled
Find all capacity reservations configured to automatically renew on their expiry date, along with the source and destination reservation IDs involved in the renewal chain.

```sql+postgres
select
  display_name,
  reservation_order_id,
  sku_name,
  term,
  expiry_date,
  renew_source,
  renew_destination
from
  azure_capacity_reservation
where
  renew = true;
```

```sql+sqlite
select
  display_name,
  reservation_order_id,
  sku_name,
  term,
  expiry_date,
  renew_source,
  renew_destination
from
  azure_capacity_reservation
where
  renew = 1;
```

### Get utilization details for capacity reservations
Analyze the utilization of your reserved VM capacity to identify underutilized reservations that may be candidates for resizing, consolidation, or cancellation.

```sql+postgres
select
  display_name,
  reservation_order_id,
  sku_name,
  quantity,
  utilization as utilization_details
from
  azure_capacity_reservation;
```

```sql+sqlite
select
  display_name,
  reservation_order_id,
  sku_name,
  quantity,
  utilization as utilization_details
from
  azure_capacity_reservation;
```

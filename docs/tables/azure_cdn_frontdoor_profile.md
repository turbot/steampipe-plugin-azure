---
title: "Steampipe Table: azure_cdn_frontdoor_profile - Query Azure CDN Front Door Profiles using SQL"
description: "Allows users to query CDN Front Door Profiles in Azure, providing detailed information about each profile, including its endpoint configurations, origin groups, and routing rules."
folder: "CDN"
---

# Table: azure_cdn_frontdoor_profile - Query Azure CDN Front Door Profiles using SQL

An Azure CDN Front Door Profile is a collection of settings and configurations for a content delivery network (CDN) front door, which is used to accelerate the delivery of web content to users globally.

## Table Usage Guide

The `azure_cdn_frontdoor_profile` table provides insights into CDN Front Door Profiles within Azure. As an Infrastructure Engineer, explore detailed information about each front door profile through this table, including its endpoint configurations, origin groups, and routing rules. Use this table to manage and optimize your CDN profiles, ensuring fast and reliable content delivery for your applications.

## Examples

### Basic front door profile information
Retrieve basic information about your Azure CDN Front Door Profiles, including their names, locations, and provisioning states.

```sql+postgres
select
  name,
  location,
  provisioning_state
from
  azure_cdn_frontdoor_profile;
```

```sql+sqlite
select
  name,
  location,
  provisioning_state
from
  azure_cdn_frontdoor_profile;
```

### SKU name and kind
Explore the SKU names and kinds of your Azure CDN Front Door Profiles, which can help in understanding the pricing tiers and types of profiles in use.

```sql+postgres
select
  name,
  sku_name,
  kind
from
  azure_cdn_frontdoor_profile;
```

```sql+sqlite
select
  name,
  sku_name,
  kind
from
  azure_cdn_frontdoor_profile;
```

### Profiles in failed provisioning state
Retrieve information about front door profiles that are in a failed provisioning state. This can help in identifying and troubleshooting issues with profile deployment.

```sql+postgres
select
  name,
  location,
  provisioning_state
from
  azure_cdn_frontdoor_profile
where
  provisioning_state = 'Failed';
```

```sql+sqlite
select
  name,
  location,
  provisioning_state
from
  azure_cdn_frontdoor_profile
where
  provisioning_state = 'Failed';
```

### Front door IDs and origin response timeout
Explore the IDs of the front doors and the origin response timeout settings of your CDN profiles.

```sql+postgres
select
  name,
  front_door_id,
  origin_response_timeout_seconds
from
  azure_cdn_frontdoor_profile;
```

```sql+sqlite
select
  name,
  front_door_id,
  origin_response_timeout_seconds
from
  azure_cdn_frontdoor_profile;
```

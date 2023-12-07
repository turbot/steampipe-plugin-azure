---
title: "Steampipe Table: azure_location - Query Azure Locations using SQL"
description: "Allows users to query Azure Locations, specifically to retrieve metadata about the different geographical locations within the Azure platform."
---

# Table: azure_location - Query Azure Locations using SQL

Azure Locations represent the geographical data centers where Azure resources are hosted. They are spread across the globe, enabling users to deploy resources near their customers to reduce latency and improve application performance. Each location is made up of one or more data centers equipped with server, storage, and networking hardware.

## Table Usage Guide

The `azure_location` table provides insights into the geographical locations within the Azure platform. As a cloud administrator or architect, explore location-specific details through this table, including name, regional display name, and longitude/latitude coordinates. Utilize it to plan your resource deployment strategy, ensuring optimal performance and compliance with data residency regulations.

## Examples

### Display name of each azure location
Explore which Azure locations are available by displaying their names. This is beneficial for understanding your geographic distribution options within the Azure platform.

```sql+postgres
select
  name,
  display_name
from
  azure_location;
```

```sql+sqlite
select
  name,
  display_name
from
  azure_location;
```


### Latitude and Longitude of the azure locations
Determine the geographical coordinates of various Azure locations. This is useful for mapping out data centers or planning for geo-redundancy.

```sql+postgres
select
  name,
  latitude,
  longitude
from
  azure_location;
```

```sql+sqlite
select
  name,
  latitude,
  longitude
from
  azure_location;
```
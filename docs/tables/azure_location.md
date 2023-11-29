---
title: "Steampipe Table: azure_location - Query Azure Locations using SQL"
description: "Allows users to query Azure Locations"
---

# Table: azure_location - Query Azure Locations using SQL

Azure Locations represent the regional presence of Azure resources. These locations are datacenters that are geographically dispersed and cater to specific geopolitical regions. They provide users with the flexibility to deploy Azure resources where they need them.

## Table Usage Guide

The 'azure_location' table provides insights into Azure Locations within Microsoft Azure. As a DevOps engineer, explore location-specific details through this table, including the name of the location, the region type, and the geographical information. Utilize it to uncover information about locations, such as those that are paired with other locations, the regions that are available for resource deployment, and the verification of geographical data. The schema presents a range of attributes of the Azure Location for your analysis, like the location name, region type, and geographical data.

## Examples

### Display name of each azure location
Explore the different Azure locations by identifying their names. This can help in understanding the distribution of your resources across different geographic regions.

```sql
select
  name,
  display_name
from
  azure_location;
```


### Latitude and Longitude of the azure locations
Discover the geographical coordinates of your Azure locations. This is useful for pinpointing the exact global positions of your resources, aiding in strategic planning and decision making.

```sql
select
  name,
  latitude,
  longitude
from
  azure_location;
```
# Table: azure_location

Azure offers the scale and data residency options you need to bring your apps closer to your users around the world.

## Examples

### Display name of each azure location

```sql
select
  name,
  display_name
from
  azure_location;
```


### Latitude and Longitude of the azure locations

```sql
select
  name,
  latitude,
  longitude
from
  azure_location;
```
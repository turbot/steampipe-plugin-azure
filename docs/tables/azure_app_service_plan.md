# Table: azure_app_service_plan

An App Service plan defines a set of compute resources for a web app to run. These compute resources are analogous to the server farm in conventional web hosting.

## Examples

### App service plan SKU info

```sql
select
  name,
  sku_family,
  sku_name,
  sku_size,
  sku_tier,
  sku_capacity
from
  azure_app_service_plan;
```


### List of Hyper-V container app service plan

```sql
select
  name,
  hyper_v,
  kind,
  region
from
  azure_app_service_plan
where
  hyper_v;
```


### List of App service plan that owns spot instances

```sql
select
  name,
  is_spot,
  kind,
  region,
  resource_group
from
  azure_app_service_plan
where
  is_spot;
```
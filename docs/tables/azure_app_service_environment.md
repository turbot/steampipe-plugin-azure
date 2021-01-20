# Table: azure_app_service_environment

The Azure App Service Environment provides a fully isolated and dedicated environment for securely running App Service apps at high scale.

## Examples

### List of app service environments which are not healthy

```sql
select
  name,
  is_healthy_environment
from
  azure_app_service_environment
where
  not is_healthy_environment;
```


### Virtual network info of each app service environment

```sql
select
  name,
  vnet_name,
  vnet_subnet_name,
  vnet_resource_group_name,
  internal_load_balancing_mode
from
  azure_app_service_environment;
```


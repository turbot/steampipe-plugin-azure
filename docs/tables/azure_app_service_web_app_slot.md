# Table: azure_app_service_web_app

When you deploy your web app on Linux, mobile back end, or API app to Azure App Service, you can use a separate deployment slot instead of the default production slot when running in the Standard, Premium, or Isolated App Service plan tier. Deployment slots are live apps with their host names. App content and configuration elements can be swapped between two deployment slots, including the production slot.

## Examples

### Basic info

```sql
select
  name,
  app_name,
  id,
  kind,
  state,
  type,
  last_modified_time_utc,
  repository_site_name,
  enabled
from
  azure_app_service_web_app_slot;
```

### List slots where the apps are enabled

```sql
select
  name,
  app_name,
  state,
  type,
  reserved,
  server_farm_id,
  target_swap_slot
  enabled
from
  azure_app_service_web_app_slot
where
  enabled;
```

### List slots that accept HTTP traffic (i.e only HTTPS is disabled)

```sql
select
  name,
  https_only,
  kind,
  region
from
  azure_app_service_web_app_slot
where
  not https_only;
```

### Host names of each slot

```sql
select
  name,
  host_names,
  kind,
  region,
  resource_group
from
  azure_app_service_web_app_slot;
```

### List enabled host names

```sql
select
  name,
  id,
  type,
  kind,
  enabled_host_names
from
  azure_app_service_web_app_slot;
```

### Get slot swap status of each slot

```sql
select
  name,
  type,
  slot_swap_status ->> 'SlotSwapStatus' as slot_swap_status,
  slot_swap_status ->> 'SourceSlotName' as source_slot_name,
  slot_swap_status ->> 'DestinationSlotName' as destination_slot_name
from
  azure_app_service_web_app_slot;
```

### Get site config details of each slot

```sql
select
  name,
  id,
  site_config ->> 'NumberOfWorkers' as number_of_workers,
  site_config ->> 'DefaultDocuments' as DefaultDocuments,
  site_config ->> 'NetFrameworkVersion' as NetFrameworkVersion,
  site_config ->> 'PhpVersion' as PhpVersion,
  site_config ->> 'PythonVersion' as PythonVersion,
  site_config ->> 'NodeVersion' as NodeVersion,
  site_config ->> 'PowerShellVersion' as PowerShellVersion,
  site_config ->> 'LinuxFxVersion' as LinuxFxVersion,
  site_config ->> 'WindowsFxVersion' as WindowsFxVersion,
  site_config ->> 'RequestTracingEnabled' as RequestTracingEnabled,
  site_config ->> 'RequestTracingExpirationTime' as RequestTracingExpirationTime,
  site_config ->> 'RemoteDebuggingEnabled' as RemoteDebuggingEnabled,
  site_config ->> 'RemoteDebuggingVersion' as RemoteDebuggingVersion,
  site_config ->> 'HTTPLoggingEnabled' as HTTPLoggingEnabled,
  site_config ->> 'DetailedErrorLoggingEnabled' as DetailedErrorLoggingEnabled,
  site_config ->> 'PublishingUsername' as PublishingUsername,
  site_config ->> 'AppSettings' as AppSettings,
  site_config ->> 'ConnectionStrings' as ConnectionStrings,
  site_config ->> 'MachineKey' as MachineKey,
  site_config ->> 'HandlerMappings' as HandlerMappings,
  site_config ->> 'DocumentRoot' as DocumentRoot
from
  azure_app_service_web_app_slot;
```

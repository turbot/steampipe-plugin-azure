# Table: azure_container_group

An Azure Container Group is a specific type of Azure Container Instances resource that allows you to group multiple containers together and run them as a single unit. A container group can contain one or more containers that are tightly coupled and need to be deployed and managed together. For example, you may have a microservices-based application that consists of multiple containers, such as a front-end container, a back-end container, and a database container. You can create an Azure Container Group to deploy and manage all these containers as a single entity.

## Examples

### Basic info

```sql
select
  name,
  id,
  provisioning_state,
  restart_policy,
  sku,
  region
from
  azure_container_group;
```

### Get encryption details of each group

```sql
select
  name,
  encryption_properties ->> 'VaultBaseURL' as vault_base_url,
  encryption_properties ->> 'KeyName' as key_name,
  encryption_properties ->> 'KeyVersion' as key_version,
  region
from
  azure_container_group;
```

### List groups that have restart policy set to `OnFailure`

```sql
select
  name,
  restart_policy,
  provisioning_state,
  type
from
  azure_container_group
where
  restart_policy = "OnFailure";
```

### Count groups by operation type

```sql
select
  os_type,
  count(name) as group_count
from
  azure_container_group
group by
  os_type;
```

### Get IP address details of each group

```sql
select
  name,
  ip_address -> 'Ports' as ports,
  ip_address ->> 'Type' as ip_address_type,
  ip_address ->> 'IP' as ip,
  ip_address ->> 'DNSNameLabel' as dns_name_label,
  ip_address ->> 'Fqdn' as fqdn
from
  azure_container_group;
```

### Get image registry credential details of each group

```sql
select
  name,
  i ->> 'Server' as server,
  i ->> 'Username' as username,
  i ->> 'Password' as password,
  i ->> 'Identity' as identity,
  i ->> 'IdentityURL' as identity_url
from
  azure_container_group,
  jsonb_array_elements(image_registry_credentials) as i;
```

### Get DNS configuration details of each group

```sql
select
  name,
  id,
  dns_config -> 'NameServers' as name_servers,
  dns_config ->> 'SearchDomains' as search_domains,
  dns_config ->> 'Options' as options
from
  azure_container_group;
```

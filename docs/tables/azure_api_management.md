# Table: azure_api_management

Azure API Management is a fully managed service that enables customers to publish, secure, transform, maintain, and monitor APIs.

## Examples

### Public and private IP address info of each API management

```sql
select
  name,
  public_ip_addresses,
  private_ip_addresses
from
  azure_api_management;
```


### API management publisher info

```sql
select
  name,
  publisher_name,
  publisher_email
from
  azure_api_management;
```


### List of premium API managements and their computing capacity

```sql
select
  name,
  sku_name,
  sku_capacity
from
  azure_api_management
where
  sku_name = 'Premium';
```


### List of API management without application tag key

```sql
select
  name,
  tags
from
  azure_api_management
where
  not tags :: JSONB ? 'application';
```

# Table: azure_ad_user

Azure Active Directory (Azure AD) is Microsoft's cloud-based identity and access management service, which helps employees sign in and access resources.

## Examples

### Basic active directory user info

```sql
select
  display_name,
  user_principal_name,
  given_name,
  mail,
  account_enabled,
  object_id
from
  azure_ad_user;
```


### List of guest users in the active directory

```sql
select
  display_name,
  user_principal_name,
  mail,
  user_type,
  usage_location
from
  azure_ad_user
where
  user_type = 'Guest';
```


### Password profile info of each user

```sql
select
  display_name,
  user_principal_name,
  additional_properties -> 'passwordProfile' -> 'enforceChangePasswordPolicy' as enforce_change_password_policy,
  additional_properties -> 'passwordProfile' -> 'forceChangePasswordNextLogin' as change_password_next_login
from
  azure_ad_user;
```


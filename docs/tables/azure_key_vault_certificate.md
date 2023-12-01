# Table: azure_key_vault_certificate

A Key Vault certificate in Azure is a managed digital certificate that provides a secure way to handle encryption keys and secrets, like passwords and .PFX files, in Azure Key Vault. Azure Key Vault is a cloud service that provides a secure store for secrets, keys, and certificates. Key Vault certificates are used in a variety of scenarios where secure communication and identity verification are needed.

## Examples

### Basic info

```sql
select
  name,
  vault_name,
  enabled,
  created,
  updated
from
  azure_key_vault_certificate;
```

### List disabled certificates

```sql
select
  name,
  vault_name,
  enabled
from
  azure_key_vault_certificate
where
  not enabled;
```

### List certificates that are expired in 10 days

```sql
select
  name,
  enabled,
  not_before,
  created,
  expires
from
  azure_key_vault_certificate
where
  expires >= now() - interval '10' day;
```

### Get key properties of certificates

```sql
select
  name,
  id,
  key_properties ->> 'Exportable' as exportable,
  key_properties ->> 'KeyType' as key_type,
  key_properties ->> 'KeySize' as key_size,
  key_properties ->> 'ReuseKey' as reuse_key
from
  azure_key_vault_certificate;
```

### Get X509 properties of certificates

```sql
select
  name,
  id,
  x509_certificate_properties ->> 'Subject' as subject,
  x509_certificate_properties -> 'Ekus' as ekus,
  x509_certificate_properties -> 'SubjectAlternativeNames' as subject_alternative_names,
  x509_certificate_properties ->> 'KeyUsage' as key_usage,
  x509_certificate_properties ->> 'ValidityInMonths' as validity_in_months
from
  azure_key_vault_certificate;
```

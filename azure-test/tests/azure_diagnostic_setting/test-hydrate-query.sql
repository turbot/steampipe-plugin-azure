select name, akas, title
from azure.azure_diagnostic_setting
where name = '{{ resourceName }}';

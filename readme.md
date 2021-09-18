# Datadog Log Exporter
To use you will need a [datadog api key](https://app.datadoghq.com/account/settings#api) and [datadog application key](https://app.datadoghq.com/access/application-keys) which can be made by following these links:
- https://app.datadoghq.com/account/settings#api
- https://app.datadoghq.com/access/application-keys

```sh
export DD_SITE="datadoghq.com" DD_API_KEY="<DD_CLIENT_API_KEY>" DD_APP_KEY="<DD_CLIENT_APP_KEY>"
go run main.go
```

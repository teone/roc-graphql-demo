# roc-graphql-demo

This is a demo implementation for a GraphQL server to interact with the Aether ROC APIs.

**Not guaranteed to work at any time :)**

## Deploy

In order to run this you need to have the `aether-roc-umbrella` chart deployed somewhere,
and you should have some data in there (ideally post the MEGA Patch).

Once that is running you need to expose the `onos-config` service with:
```shell
kubectl port-forward --address 0.0.0.0 svc/onos-config 5150
```
_This command assumes everything is deployed in the `default` namespace._

Then you can run this server with:
```shell
go run cmd/server.go
```

And at that point you can open the GraphQL server at `http://localhost:8080`
# sb-tool
tool to working with azure service bus in go

Example config file

```
ConnectionString: "insert your connection string"

Producer:
  Topic: "topic"
  Upload:
    Interval: "10s"
    Bytes: 16000
    Timeout: "5s"

Consumer:
  Topic: "topic"
  Subscription: "subscription"
  Timeout: "1m"
```

To run 

```
go run main.go producer --config ./config.yaml
go run main.go consumer --config ./config.yaml 
```

Producer uploads messages at the rate and msg payload size specified

Consumer reads messages off of a service bus as fast as it can.
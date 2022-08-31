#### Hexa RPC is Hexa-related RPC & gRPC SDK

#### Requirements:
go : minimum version `1.13`

#### Install
```
go get github.com/kamva/hexa-rpc
```

### Proposal
- [ ] Change Hexa error status from http status to gRPC status, because:    

  __Advantages__  
  - If we using gRPC status, so will dont need to convert https status to gRPC status, we will just convert our gRPC status to http on the gateway responses.   

  __Drawbacks__  
  - If we use gRPC codes, so we will need to import gRPC library in all of libraries that need to define or use hexa error, while http statues exists in most languages  as standard libraries.
  

#### Todo
- [ ] Use `recover` interceptor in the [gRPC interceptors](https://github.com/grpc-ecosystem/go-grpc-middleware).
- [x] Implement status to Hexa error (and reverse) mapper.
- [x] Set Hexa logger as gRPC Logger (implement gRPC logger adapter by hexa logger)
- [x] Implement request logger (log start-time, end-time, method, error,...)
- [ ] We should implement all of our interceptors for the Stream request/responses also (for now
 we just support Unary Request/responses).
- [ ] Write Tests
- [ ] Add badges to README.
- [ ] CI

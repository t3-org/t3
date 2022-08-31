#### hecho (hexa echo) contain middlewares,handlers,... for the echo.

#### Install
```
go get github.com/kamva/hecho
```

##### Middlewares
* log: set new log handler as context log that contains:
    - request id in eac log record.
    - users data in each log record.

* transltion: Set new translator in context that localized with
users accept-languages and then fallback and default languages.


##### Handlers
* error handler: handle hexa errors.
    
##### Middleware dependencies:
* `hecho.CurrentUser` middleware requires
    - `hecho.JWT` middleware (load `JWT` middleware before `CurrentUser`).
* `hecho.HexaContext` middleware requires 
    - echo `middleware.RequestID`
    - hexa `hecho.CorrelationID`
    - hexa `hecho.CurrentUser` middleware.
* `hecho.SetContextLogger` middleware requires
    - hexa `hexa.HexaContext`
* `hecho.TracingDataFromUserContext` middleware requires
  - hexa `hecho.HexaContext`

#### Todo:
- [ ] Map echo errors (see errors list in `echo.go:263`) to hexa error with translation.
- [ ] Tests
- [ ] Add badges to readme.
- [ ] CI 


### Hexa arranger is a wrapper and helper for Uber Cadence

### Install
```bash
go get github.com/kamva/hexa-arranger
```


##### Notes about error handling:
- Read [the docs](https://github.com/temporalio/sdk-go/blob/master/temporal/error.go) about error types in temporal.

- solution to handle errors:  
    temporal doesn't have interceptor for activities yet.
    For workflows but it has, maybe we use it later.
    It should report and convert errors when return 
    error from either activity or workflow.
    I think a simple decorator for each workflow
    and activity can be a good idea. when temporal
    implemented interceptor, so we can replace it with an interceptor.

- note: if you got error from activity in your workflow and want to return it, do not use `tracer.Trace` as wrapper of that error, because if error is temporal error, and we wrap it using another error, so temporal worker will get custom error and try to convert it to ApplicationError, so original activity error will
 wrap by another ApplicationError and if first called activity converted hexa error to ApplicationError, we will 
 do not get it because that ApplicationError is wrapped by another new ApplicationError now.
  
error behaviour:

```
-> when activity handler returns error:
    - report error
    - convert hexa error to application error

-> when workflow get activity|child_workflow error:
    - nothing. it can convert it to hexa error if needed(using our implmented error converter funcitons).

-> when workflow itself returns error:
    - report error 
    - convert hexa error to application error
    
-> When our app get returned workflow error:
      - convert error to hexa (simply by our implmemented error converters).
      - report error (don't need to implement, our app report every error in its edge like router or gRPC interceptors,...)
```


#### TODO
- [ ] Provide OpenTracing.
- [ ] Write tests.

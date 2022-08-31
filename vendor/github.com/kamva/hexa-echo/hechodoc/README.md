Conventions:
---
- You must create openapi go file initially.
  (e.g., `openapi_docs.go`)
  

Notes:
-------
- `raw route name` is the raw route name which we set on route definition in 
  our source code. e.g., `hi::mehran`
- `route name` in the source code is camelcased value of the
  raw route name e.g., `hiMehran`.
  
  
How to generate docs?
---
- Install swagger
- Place a package named something like `doc`.
- In the `doc` package create a file named something like `openapi_docs.go`.
  This is your openapi comments file, so all of your
  documentation should be in this file.
- In the `doc` package create another package named `gen`
- In your main file import `doc` package as blank. (e.g., `_ "github.com/kamva/hexa-echo/examples/docexample/doc"`)
- Now in you source code create two commands: `extract`,`trim`.
- In order run `extract template` and `trim template` commands.
  sample in the `examples` package´.

FAQ:
---
- How should I install swagger command? simply install
it as go package: 
```bash
GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger
```


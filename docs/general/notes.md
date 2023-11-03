### DB improvement
- Create two structs (`Table`, `Col`) to return Identifiers of sea_query library.
- Read more about custom types [here](/Users/mehran/.cargo/registry/src/index.crates.io-6f17d22bba15001f/sqlx-0.7.2/src/macros/mod.rs)
  to see when we should do:
   - `as "!: Option<Level>"` for our custom types in selection
   - use also `as Option<Level>` when we set/insert a custom typed field when 
     using `query!` macro instead of runtime call.
- Add a method to set custom error on db constraints(e.g., when a row is duplicate)
see [here](https://github.com/launchbadge/realworld-axum-sqlx/blob/main/src/http/error.rs#L199)


### architecture improvement
- let services have `Option<T>` as their needed servres and also an `inline functino` for 
  each service to return the services without calling to `unwrap` method.


### Matrix Links
- [matrix concepts](https://spec.matrix.org/v1.8/client-server-api/#sending-events-to-a-room)
- [matrix api](https://spec.matrix.org/v1.8/client-server-api/#sending-events-to-a-room)
- [matrix sdk examples](https://github.com/matrix-org/matrix-rust-sdk/tree/main/examples)
- [matrix element client](https://app.element.io)
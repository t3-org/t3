Use [hashicorp go-plugin](https://github.com/hashicorp/go-plugin) package to implement the plugin system.

Usage:

- Add sources
- Add channels
- Add other capabilities like AI assistant.

__Notes:__
- Provide an `API` interface to let plugins communicate with our app.
- On startup of the plugin give it the `API` instance to let them keep the `API` and use it when needed. (e.g., after client creation of the plugin, call a method like `Init` with `API` as argument).
- Provide a `ServeHTTP()` method on for each plugin to serve http requests after installation (provide bodyReader and responseWriter through the go-plugin package).

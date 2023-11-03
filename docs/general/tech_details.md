To implement channels:

```rust

// Have an enum of all channel types.
pub enum ChannelType {
    Matrix(MatrixChannel),
    // Add other channel types here.
}

pub struct MatrixChannel {} // implement the channel itself. get the `app: Ac<App>` at creation time.


// Each channel must implement this interface:
pub trait Channel {
    fn start();
    fn firing(ticket: Ticket);
    fn resolved(ticket: Ticket);
}


// Define a config to let use define all channels.
// here we should iterate over all defined user's channels and create new channel instance by the specified
// channel type and append it to the channels map. the map's type is HashMap<String,ChannelType>.


// In the app layer find the ChannelType of the webhook(using the provider channels' map) and its `firing` 
// or `resolved` method to inform its target third-party service.

// At the app startup time per each channel of the channel map in the provider spawn a task and  call its `start`
// method to start listening to the channel's incoming messages.
```
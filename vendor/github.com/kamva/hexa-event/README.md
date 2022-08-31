#### hexa event implement events feature for Hexa

#### Install

```
go get github.com/kamva/hexa-event
```

__Supported message queues and event streaming platforms :__

- Kafka & Kafka outbox pattern (supports deduplication)
- Pulsar
- Nats streaming

__Known errors :__

- when two time subscribe with same subscription name with subscription type : `pulsar.Exclusive` on __pulsar__ driver,
  you get this error :

 ```bash
server error: ConsumerBusy: Exclusive consumer is already connected
```

#### Proposal:

- [ ] Remove the `HandlerContext` as first param and get error as return param of `EventHandler`, if you got an error,
  so return negative signal and log the error, otherwise return positive signal to the event broker.
- [ ] Remove the `err` param as lat para of each handler, if occured error, so just log it and send nack, because we
  should not get any error in our app, if we get error on an event, so we don't need to call to the handler, we need to
  fix it.
- [ ] Transactional outbox: we need to implement another emitter for each each driver which instead of emit, it prepare
  the message to emit later. we provide the hexa Context with its propagators, maybe we need to provide some data with
  that event which is avaiable in that time like hexa Context, we should convert message to raw message and do what we
  need to do to send a message, then store it in that driver's store.
- [ ] After that we also need to another emitter interface which gets each document and emit that event directory. we
  also have another option to recreate event and emit it using noraml emiter, but this is not good, we don't need to
  recreate the event, also we have some data like propaged hexa Context which we realy don't need to reacreate it to
  emit event, some emittters like pulsar emitter add some data on sending event like "send_time", we should implement
  another emitter for db events, which does not use the normal emitter interface, because it need to get the events. we
  use this or some other thing which gets an struct which return db doc and on get new doc, it should emit it.

#### Todo:

- [x] Add support of protocol buffer to send events.
- [x] Add `Extra []interface{}` option to the `SubscriptionOptions` to get more features on each subscription relative
  to each driver. remove list of options in consumerGenerator(we can generate without a consumer generator or simple
  consumer generator) __[Accepted]__.
- [x] Implement nats-streaming driver
- [ ] Implement a new background process to remove old messages in the kafka outbox pattern.
- [ ] Implement Mock driver
- [ ] Write Tests.
- [ ] Implement mock
- [ ] Add badges to readme.
- [ ] CI.

#### Implementation notes:

- If you need to handle message deduplication, you should add following keys to your hexa context:
    - `HEXA_EVENT_ID`: A unique ID for the current event.
    - `HEXA_EVENT_HANDLER_ACTION_NAME`: The action's name that your event handler want to do. we use this value to do
      not conflict deduplication when you listen to one event multiple times for multiple actions in one microservice.
    - `HEXA_ROOT_EVENT_ID`: Its id of the base event in retry events (the root event's id). for the root event itself is
      equal to the eventId.
    - `HEXA_ROOT_EVENT_HANDLER_ACTION_NAME`: The action's name of the root event (the base event of the retry events)
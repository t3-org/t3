package hevent

// SubscribeMulti subscribes to multiple channel.
func SubscribeMulti(r Receiver, options ...*SubscriptionOptions) error {
	for _, o := range options {
		if err := r.SubscribeWithOptions(o); err != nil {
			return err
		}
	}

	return nil
}



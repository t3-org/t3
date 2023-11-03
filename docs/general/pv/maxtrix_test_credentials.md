My credentials for my test account to be used in my matrix tests:

homeserver: https://matrix-client.matrix.org
username: mehran-bot
password: n5Vy2]dyXU4Hri3

You can browse this account on the element page in the firefox(I've logged in firefox): https://app.element.io

Example(in the `playground/marix-sdk-rust`
dir): `cargo run -p example-getting-started https://matrix-client.matrix.org mehran-bot n5Vy2]dyXU4Hri3`
And then open "https://app.element.io" on chrome(which contains my real account, not this test account) and then
send `!party` on the `webhook_check` room to see
the response that we'll get from the app that we've run.


### Links

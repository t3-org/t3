[got from here](https://blog.facha.dev/how-to-self-host-matrix-and-element-docker-compose/)

- run it: `docker compose up -d`
- Create two users with the following specs :
    - username: admin, password: admin, isAdmin: true
    - username: bot, password: botpass, isAdmin: false

  Use `register_new_matrix_user -c /etc/homeserver.yaml http://localhost:8008` command
  in the synapse container, to create the user.

- Open `http://localhost:8080`, change homeserver to "http://localhost:8008" and then login.
- Create a room, called `t3` from element ui.
- invite `@bot:matrix.example.com` into the `t3` room.
- Done.

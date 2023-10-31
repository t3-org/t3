CREATE TABLE ticket_values
(
    ticket_id bigint  not null references tickets (id) on delete cascade,
    key       text NOT NULL,
    value     text  NOT NULL,
    UNIQUE (ticket_id, key)
)

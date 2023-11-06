create type TicketLevel AS ENUM ('low','medium','high');
CREATE TABLE tickets
(
    id          bigserial PRIMARY KEY,

    annotations text,
    fingerprint text UNIQUE NOT NULL,
    is_firing   BOOLEAN     not null,
    started_at  bigint      NOT NULL,
    ended_at    bigint,

    is_spam     boolean     not null,
    level       TicketLevel,
    description text,
    seen_at     bigint,
    created_at  bigint,
    updated_at  bigint
);

-- ticket indexes
create index on tickets (fingerprint);
create index on tickets (is_firing);
create index on tickets (is_spam);
create index on tickets (level);
create index on tickets (seen_at);
create index on tickets (started_at, ended_at);

create table ticket_labels
(
    ticket_id bigint not null references tickets (id) on delete cascade,
    key       text   not null,
    val       text   not null,
    UNIQUE (ticket_id, key)
);


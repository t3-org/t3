CREATE TABLE IF NOT EXISTS planets
(
    id         VARCHAR(255) primary key,
    name       varchar(255) unique,
    code       VARCHAR(255),
    updated_at bigint,
    created_at bigint
);

create index if not exists idx_planets_created_at on planets (created_at);
create index if not exists idx_planets_updated_at on planets (updated_at);
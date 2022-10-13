CREATE TABLE IF NOT EXISTS systems
(
    name       VARCHAR(64) PRIMARY KEY,
    value      VARCHAR(1024),
    updated_at bigint,
    created_at bigint
);

create index if not exists idx_systems_created_at on systems (created_at);
create index if not exists idx_systems_updated_at on systems (updated_at);

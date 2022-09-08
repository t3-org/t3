CREATE TABLE IF NOT EXISTS planets
(
    name       VARCHAR(255) primary key,
    code      VARCHAR(255),
    updated_at timestamptz,
    created_at timestamptz
);

create index if not exists idx_planets_created_at on planets (created_at);
create index if not exists idx_planets_updated_at on planets (updated_at);
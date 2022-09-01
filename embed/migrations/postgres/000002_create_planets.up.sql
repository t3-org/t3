CREATE TABLE IF NOT EXISTS planets
(
    id         VARCHAR(36) PRIMARY KEY,
    name       VARCHAR(255),
    code       VARCHAR(255) unique,
    updated_at timestamptz,
    created_at timestamptz
);

create index if not exists idx_planets_created_at on planets (created_at);
create index if not exists idx_planets_updated_at on planets (updated_at);




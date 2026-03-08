create table if not exists blog (
    id serial primary key,
    created_at timestamptz not null default now(),
    modified_at timestamptz not null default now(),
    is_deleted bool not null default false,
    content text
);

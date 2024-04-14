create table if not exists banners (
    id serial,
    tag_ids integer[],
    feature_id integer,
    primary key (tag_ids, feature_id),
    content json,
    is_active bool,
    created_at timestamp,
    updated_at timestamp
);

create table if not exists tags (
    id serial primary key,
    name varchar(255)
);

create table if not exists features (
    id serial primary key,
    name varchar(255)
);
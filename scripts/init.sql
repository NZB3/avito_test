create table if not exists banners (
    id serial,
    tag_ids integer[],
    feature_id integer,
    primary key (tag_ids, feature_id),
    content json,
    is_active bool,
    created_at timestamp default now(),
    updated_at timestamp default now()
);

create or replace function set_updated_at() returns trigger as
$$
    BEGIN
        NEW.updated_at := current_timestamp;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

create or replace trigger update_trigger before update on banners
    for statement execute procedure set_updated_at();

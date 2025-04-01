-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE order_status AS ENUM ('active', 'canceled');

create table if not exists orders (
                                      id uuid not null primary key,
    user_id uuid not null,
                                      status order_status default 'active',
                                      created_at TIMESTAMP default NOW(),
                                      updated_at TIMESTAMP default Now(),
    constraint fk_constraint foreign key (user_id)
                                  references users(id)
);

create table if not exists users (
                                     id uuid not null primary key,
                                     first_name varchar(255) not null,
                                     last_name varchar(255) not null,
                                     email varchar(255) not null unique,
                                     phone varchar(10) not null unique,
                                     created_at TIMESTAMP default NOW(),
                                     updated_at TIMESTAMP default Now(),
    order_history uuid[] default '{}'::uuid[]
);

create table if not exists inventories (
    id serial primary key,
    product int not null,
    quantity int default 0,
    price float default 0.0,
    created_at TIMESTAMP default NOW(),
    updated_at TIMESTAMP default Now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table inventories;
drop table orders;
drop table users;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;

create table if not exists users (
                                     id uuid not null primary key,
                                     first_name varchar(255) not null,
                                     last_name varchar(255) not null,
                                     email varchar(255) not null unique,
                                     phone varchar(10) not null unique,
                                     created_at TIMESTAMP default NOW(),
                                     updated_at TIMESTAMP default Now()
);

create table if not exists categories (
    id serial primary key,
    name varchar(255) not null unique,
    created_at TIMESTAMP default NOW(),
    updated_at TIMESTAMP default Now()
);

create table if not exists products(
    id serial primary key,
    sku varchar(255) unique not null,
    name varchar(255) not null unique,
    available bool default false,
    created_at TIMESTAMP default NOW(),
    updated_at TIMESTAMP default Now(),
    price float not null,
    reviews varchar[] default '{}'::varchar[],
    description text not null,
    image varchar(255) not null,
    category int not null,
    CONSTRAINT fk_constraint FOREIGN KEY (category)
                                   references categories(id) on delete cascade
);

create table if not exists carts (
    id serial primary key,
    items int[] default '{}'::int[],
    owner uuid not null,
    created_at TIMESTAMP default NOW(),
    updated_at TIMESTAMP default Now(),
    constraint fk_constraint foreign key (owner)
                                 references users(id) on delete cascade
);

create table if not exists inventories (
    id serial primary key,
    product int not null,
    quantity int default 0,
    created_at TIMESTAMP default NOW(),
    updated_at TIMESTAMP default Now(),
    constraint fk_constraint foreign key (product)
                                       references  products(id) on delete cascade
);

CREATE TYPE order_status AS ENUM ('active', 'canceled');

create table if not exists orders (
    id uuid not null primary key,
    status order_status default 'active',
    created_at TIMESTAMP default NOW(),
    updated_at TIMESTAMP default Now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
drop table carts;
drop table products;
drop table inventories;
drop table categories;
drop table orders;
-- +goose StatementEnd

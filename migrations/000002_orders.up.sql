create table orders
(
    id         uuid        not null
        constraint orders_pk
            primary key,
    courier_id uuid        null,
    location_x integer     not null,
    location_y integer     not null,
    volume     integer     not null,
    status     varchar(32) not null
);
create table couriers
(
    id         uuid
        constraint couriers_pk
            primary key,
    name       varchar(255) not null,
    speed      int          not null,
    location_x int          not null,
    location_y int          not null
);

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

create table storage_places
(
    id         uuid         not null
        constraint storage_place_pk
            primary key,
    name       varchar(255) not null,
    volume     integer      not null,
    order_id   uuid         null,
    courier_id uuid         null
);

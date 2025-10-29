create table storage_places
(
    id         uuid         not null
        constraint storage_place_pk
            primary key,
    name       varchar(255) not null,
    volume     integer      not null,
    order_id   uuid         null,
    courier_id uuid         not null
);
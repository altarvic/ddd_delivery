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
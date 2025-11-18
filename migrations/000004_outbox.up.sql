create table outbox
(
    id        uuid                     not null
        constraint outbox_pk
            primary key,
    type      TEXT                     not null,
    content   TEXT                     not null,
    occurred  TIMESTAMP with time zone not null,
    processed TIMESTAMP with time zone
);
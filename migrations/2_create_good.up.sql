create table good
(
    id         serial
        constraint good_pk
            primary key,
    name       varchar(255) not null,
    owner_id   int          not null
        constraint good_trader_id_fk
            references trader
            deferrable initially deferred,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);

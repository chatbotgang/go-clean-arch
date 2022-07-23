create table trader
(
    id         serial
        constraint trader_pk
            primary key,
    uid        varchar(36)  not null,
    email      varchar(255) not null,
    name       varchar(36)  not null,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);

create unique index trader_uid_uniq
    on trader (uid);

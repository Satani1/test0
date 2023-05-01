create table test
(
    order_uid varchar not null,
    data      jsonb,
    ids       serial
        constraint test_pk
            primary key
);

alter table test
    owner to postgres;


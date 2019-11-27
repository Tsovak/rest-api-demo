create table if not exists accounts (
    id bigserial not null primary key,
    name varchar(256) not null,
    currency varchar(32) not null,
    balance int not null check (balance >= 0),
    unique (id)
);

create table if not exists payments (
    id bigserial not null primary key ,
    amount bigserial not null,
    to_account_id varchar(256) not null,
    from_account_id varchar(256) not null,
    direction varchar(256) not null
);

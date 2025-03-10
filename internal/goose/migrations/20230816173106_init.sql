-- +goose Up
-- SQL in this section is executed when the migration is applied.
create schema slyip;

create table slyip.account
(
    id                   uuid primary key         not null default gen_random_uuid(),
    first_name           varchar(255)             not null default '',
    last_name            varchar(255)             not null default '',
    phone                varchar(255)             not null default '',
    email                varchar(255)                      default null unique,
    is_email_verified    boolean                  not null default false,
    is_phone_verified    boolean                  not null default false,
    password_hashed      varchar(255)             not null default '',
    invitation_code      varchar(255)             not null default '',
    role                 varchar(255)             not null default 'basic',
    last_used_sly_wallet varchar(255)             not null default '',
    created_at           timestamp with time zone not null default now(),
    updated_at           timestamp with time zone not null default now()
);

create table slyip.sly_wallet
(
    address            varchar(255) primary key not null unique,
    chainId            varchar(255)             not null,
    account_id         uuid                     not null,
    transaction_hash   varchar(255)             not null default '',
    transaction_status varchar(255)             not null default 'not_initiated',
    invitation_code    varchar(255)             not null,
    constraint FK_acc_wallet foreign key (account_id) references slyip.account (id),
    created_at         timestamp with time zone not null default now(),
    updated_at         timestamp with time zone not null default now()
);

create table slyip.ecdsa
(
    address    varchar(255) primary key not null unique,
    account_id uuid                     not null,
    constraint FK_acc foreign key (account_id) references slyip.account (id),
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

create table slyip.ecdsa_sly_wallet
(
    ecdsa_address            varchar(255) not null,
    on_chain_account_address varchar(255) not null,
    on_chain_permissions     int          not null default 0,
    primary key (ecdsa_address, on_chain_account_address),
    constraint FK_ecdsa foreign key (ecdsa_address) references slyip.ecdsa (address),
    constraint FK_sly foreign key (on_chain_account_address) references slyip.sly_wallet (address)
);


create table slyip.invitation_code
(
    code             varchar(255) primary key not null,
    transaction_hash varchar(255)             not null default '',
    expires_at       timestamp with time zone not null
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
-- SQL in this section is executed when the migration is applied

drop table slyip.invitation_code;
drop table slyip.ecdsa_sly_wallet;
drop table slyip.ecdsa;
drop table slyip.sly_wallet;
drop table slyip.account;
drop schema slyip;



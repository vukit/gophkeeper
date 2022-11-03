create table users (
    "user_id"   serial primary key,
    "username"  varchar(64) not null,
    "password"  char(64) not null,
    unique ("username")
);

create index "users_username_idx" ON users ("username");

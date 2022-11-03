create table logins (
    "login_id"  serial primary key,
    "user_id"   int not null references users on delete cascade,
    "username"  character varying not null,
    "password"  character varying not null,
    "metainfo"  character varying not null
);

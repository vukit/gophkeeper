create table files (
    "file_id"   serial primary key,
    "user_id"   int not null references users on delete cascade,
    "path"      character varying not null,
    "name"      character varying not null,
    "metainfo"  character varying not null
);

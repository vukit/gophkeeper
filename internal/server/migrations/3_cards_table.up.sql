create table cards (
    "card_id"   serial primary key,
    "user_id"   int not null references users on delete cascade,
    "bank"      character varying not null,
    "number"    character varying not null,
    "date"      character varying not null,
    "cvv"       character varying not null,
    "metainfo"  character varying not null
);

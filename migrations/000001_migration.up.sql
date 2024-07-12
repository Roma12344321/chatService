create table if not exists person
(
    id       int generated by default as identity primary key,
    username varchar(100) not null unique,
    password varchar(100) not null,
    role     varchar(30)  not null default 'ROLE_USER'
);

create table if not exists chat_room
(
    id   int generated by default as identity primary key,
    name varchar(100) not null,
    date timestamp    not null
);

create table if not exists person_chat_room
(
    person_id    int references person (id) on delete cascade,
    chat_room_id int references chat_room (id) on delete cascade,
    primary key (person_id, chat_room_id)
);

create table if not exists message(
    id bigint generated by default as identity primary key,
    text varchar(150),
    person_id int references person(id) on delete set null,
    chat_room_id int references chat_room(id) on delete cascade
);
create table if not exists users (
    id int generated always as identity primary key, 
    email varchar(254) not null,
    name text not null,
    role varchar(254) not null, 
    created_at timestamp with time zone not null
);

create table if not exists expenses (
    id int generated always as identity primary key,
    user_id int references users(id) on delete restrict, 
    amount_idr bigint not null,
    description text not null,
    receipt_url text not null,
    status varchar(50) not null,
    submitted_at timestamp with time zone not null,
    processed_at timestamp with time zone not null
)
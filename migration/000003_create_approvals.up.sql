create table if not exists approvals (
    id int generated always as identity primary key,
    expense_id int not null references expenses(id) ON DELETE NO ACTION,
    approver_id int references users(id)  ON DELETE NO ACTION,
    status varchar(50) not null,
    notes text not null,
    created_at timestamp with time zone
)
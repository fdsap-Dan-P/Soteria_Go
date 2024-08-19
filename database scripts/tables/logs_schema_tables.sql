-- audit tables
create table logs.audit_logs (
        user_id int not null,
        staff_id text not null,
        username text not null,
        user_status text not null,
        institution_name text not null,
        user_action text not null,
        old_value text null,
        new_value text null,
        log_date text not null,
        log_time text not null,
        created_at timestamp default current_timestamp,
        updated_at timestamp default current_timestamp
);


-- for user log
create table logs.login_logs (
        log_id serial primary key,
        user_id INT references user_accounts(user_id) on delete cascade,
        ip text,
        login_attempt_at text,
        is_success bool,
        created_at timestamp default current_timestamp,
        updated_at timestamp default current_timestamp
);


CREATE TABLE logs.locked_users (
    user_id INTEGER,
    ip TEXT,
    created_at TEXT
);
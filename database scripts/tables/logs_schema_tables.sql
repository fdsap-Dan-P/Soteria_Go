CREATE TABLE logs.user_tokens (
        token_id serial primary key,
        username TEXT,
        staff_id TEXT, 
        token TEXT,
        insti_code TEXT references offices_mapping.institutions(institution_code) on delete cascade,
        app_code TEXT references public.applications(app_code) on delete cascade,
        created_at timestamp default current_timestamp,
);
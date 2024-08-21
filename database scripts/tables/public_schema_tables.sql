-- tables for users 
create table public.user_accounts (
	user_id serial primary key,
	username text unique not null,
	first_name text not null,
	middle_name text null,
	last_name text not null,
	email text unique not null,
	phone_no text not null,
	staff_id text not null,
	last_login text null,
    institution_id int references offices_mapping.institutions(institution_id) on delete cascade,
	created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP
);

create table public.user_passwords (
	user_id int references public.user_accounts(user_id) on delete cascade,
	password_hash bytea NULL,
	requires_password_reset bool null default true,
	last_password_reset text null,
	created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP
);

-- tables for applications 
create table public.applications (
	applicayion_id serial primary key,
	application_code text not null,
	application_name text not null,
	application_description text null,
	api_key text null, --auto generation upon registration
	created_at timestamp default current_timestamp,
	updated_at timestamp default current_timestamp
);

insert into public.applications
	(app_code, app_name, app_description) 
	values
	('data-platform', 'Data Platform', 'Data Platform Web Application for Reports'),
	('cagabay', 'Cagabay', 'Cagabay Mobile Application')
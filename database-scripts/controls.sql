create schema parameters;
create schema logs;
create schema offices_mapping;

create user cagabay_ua_user with password 'S0t3rI@-G09*7SUf';

GRANT all privileges on database cagabay_ua_prod TO cagabay_ua_user;

GRANT usage on schema public TO cagabay_ua_user;
GRANT usage on schema parameters TO cagabay_ua_user;
GRANT usage on schema logs TO cagabay_ua_user;x
GRANT usage on schema offices_mapping TO cagabay_ua_user;

GRANT all privileges on all tables in schema public TO cagabay_ua_user;
GRANT all privileges on all tables in schema parameters to cagabay_ua_user;
GRANT all privileges on all tables in schema logs to cagabay_ua_user;
GRANT all privileges on all tables in schema offices_mapping to cagabay_ua_user;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO cagabay_ua_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA parameters TO cagabay_ua_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA logs TO cagabay_ua_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA offices_mapping TO cagabay_ua_user;

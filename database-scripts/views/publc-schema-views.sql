-- get user's details upon login
create or replace view public.user_details as
select
    ua.username,
    ua.first_name,
    ua.middle_name,
    ua.last_name,
    ua.email,
    ua.phone_no,
    ua.staff_id,
    ua.last_login,
    i.institution_code,
    i.institution_name
from
    public.user_accounts ua
    join offices_mapping.institutions i on i.institution_id = ua.institution_id
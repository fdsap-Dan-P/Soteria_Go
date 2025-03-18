-- get user's details upon login
CREATE OR REPLACE VIEW public.user_details
AS SELECT ua.user_id,
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
   FROM user_accounts ua
     JOIN offices_mapping.institutions i ON i.institution_id = ua.institution_id
     join applications a on a.application_id = ua.application_id;

-- get user's linked applications
create or replace view user_app_view as 
select 
	ua.username, 
	ua.staff_id,
	ua.email,
	ua.phone_no,
	a.application_code,
	a.application_name,
	a.application_description,
	i.institution_code,
	i.institution_name
from user_accounts ua 
inner join user_applications uapp on ua.user_id = uapp.user_id
inner join applications a on uapp.application_id = a.application_id
inner join offices_mapping.institutions i on ua.institution_id = i.institution_id;
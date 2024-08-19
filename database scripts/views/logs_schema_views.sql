-- get user details for audit trail
CREATE OR REPLACE VIEW logs.user_details_audit_trail AS
SELECT
    ua.user_id, ua.staff_id, ua.username, us.status_name, i.institution_name
FROM
    public.user_accounts ua
    LEFT JOIN public.user_status us ON ua.status_id = us.status_id
    LEFT JOIN offices_mapping.institutions i ON ua.institution_id = i.institution_id
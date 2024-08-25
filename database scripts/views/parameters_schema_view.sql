create or replace view parameters.system_config_parameters as
select
    sc.config_id,
    sc.config_code,
    sc.config_name,
    sc.config_description,
    sc.config_type,
    iac.config_value,
    iac.config_insti_code,
    iac.config_app_code,
    iac.created_at,
    iac.updated_at
from
    parameters.system_config sc
    join parameters.insti_app_config iac on iac.config_id = sc.config_id

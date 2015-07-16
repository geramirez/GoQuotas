CREATE VIEW quotas_view AS
    SELECT *, (SELECT json_agg(t) AS data FROM
        (SELECT memory, Count(date) as days FROM quotadata WHERE quotadata.guid=quotas.guid GROUP BY memory) t)
    FROM quotas;

CREATE VIEW quota_details AS
    SELECT guid, (SELECT json_agg(t) AS details FROM
        (SELECT guid, name, data FROM quotas_view WHERE quotas_view.guid=quotas.guid ) t)
    FROM quotas;


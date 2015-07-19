CREATE OR REPLACE FUNCTION get_quotas_details(parm1 DATE, param2 DATE)
  RETURNS TABLE (guid text, name text, cost numeric, data json)
  AS
    $body$
      with cost_table as (
        SELECT guid, SUM(cost) AS cost
        FROM (
          SELECT guid, memory * .0033 as cost, date
          FROM quotadata
          WHERE date between $1 and $2
        ) AS costs
        GROUP BY guid
      )

      SELECT quotas.guid, quotas.name, cost_table.cost, (
        SELECT json_agg(t) AS data
        FROM (
          SELECT memory, Count(date) as days
          FROM quotadata
          WHERE quotadata.guid=quotas.guid and date between $1 and $2
          GROUP BY memory)
      t)
      FROM quotas, cost_table
      WHERE quotas.guid = cost_table.guid;
    $body$
language sql;


CREATE OR REPLACE FUNCTION get_quotas(parm1 DATE, param2 DATE)
  RETURNS TABLE (guid text, details json)
  AS
    $body$
  with quotas_cost_view as (
    select *
    from get_quotas_details($1, $2)
  )
  SELECT guid, (
    SELECT json_agg(t) AS details
    FROM (
      SELECT guid, name, cost, data
      FROM quotas_cost_view
      WHERE quotas_cost_view.guid=quotas.guid )
    t)
  FROM quotas;
  $body$
language sql;


/*
-- quotas/:guid endpoint
select * from get_quotas('2015-01-01', '2016-01-01')
*/

/*
-- quotas/ endpoint
with quota_details_agg as (
  select * from get_quotas_details('2015-01-01', '2016-01-01')
)
SELECT json_agg(t) AS elements FROM (
  SELECT guid, name, cost, data
  FROM quota_details_agg
) t
*/

/*
-- old views
CREATE OR REPLACE  VIEW quotas_view AS
    SELECT *, (SELECT json_agg(t) AS data FROM
        (SELECT memory, Count(date) as days FROM quotadata WHERE quotadata.guid=quotas.guid GROUP BY memory) t)
    FROM quotas;

CREATE OR REPLACE VIEW quota_details AS
    SELECT guid, (SELECT json_agg(t) AS details FROM
        (SELECT guid, name, data FROM quotas_view WHERE quotas_view.guid=quotas.guid ) t)
    FROM quotas;

*/

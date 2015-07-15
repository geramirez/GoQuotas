-- psql -U go_quotas -d go_quotas -a -f src/db_init.sql
CREATE TABLE IF NOT EXISTS quotas
(
    guid VARCHAR(500) PRIMARY KEY, 
    name VARCHAR(500) NOT NULL
);

CREATE TABLE IF NOT EXISTS quotadata
(
    guid VARCHAR(500) REFERENCES quotas, 
    memory INT,
    date DATE,
    CONSTRAINT guid_date PRIMARY KEY(guid, date)
)

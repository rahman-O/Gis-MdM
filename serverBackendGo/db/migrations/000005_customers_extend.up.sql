-- Customer columns for super-admin search/CRUD parity (subset of Java Liquibase).

ALTER TABLE customers ADD COLUMN IF NOT EXISTS email VARCHAR(50);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS accounttype INT NOT NULL DEFAULT 0;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS customerstatus VARCHAR(100);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS registrationtime BIGINT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS expirytime BIGINT;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS devicelimit INT NOT NULL DEFAULT 3;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS deviceconfigurationid INT;

-- Demo tenant for control-panel impersonate smoke (non-master).
INSERT INTO customers (name, description, master, prefix, email, registrationtime)
SELECT 'Demo Tenant', 'Dev smoke tenant', FALSE, 'demo-', 'demo@localhost', EXTRACT(EPOCH FROM NOW())::BIGINT * 1000
WHERE NOT EXISTS (SELECT 1 FROM customers WHERE master = FALSE AND name = 'Demo Tenant');

INSERT INTO users (login, email, name, password, customerid, userroleid, authtoken)
SELECT
    'demo_tenant',
    'demo@localhost',
    'Demo Admin',
    '349242D38ED8667B5C11D2412EBEA4636BD3CA3A',
    c.id,
    2,
    'smokeorgadmintoken00001'
FROM customers c
WHERE c.name = 'Demo Tenant' AND c.master = FALSE
  AND NOT EXISTS (SELECT 1 FROM users u WHERE u.customerid = c.id AND u.userroleid = 2);

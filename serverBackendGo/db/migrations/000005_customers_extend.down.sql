DELETE FROM users WHERE login = 'demo_tenant';
DELETE FROM customers WHERE name = 'Demo Tenant' AND master = FALSE;

ALTER TABLE customers DROP COLUMN IF EXISTS deviceconfigurationid;
ALTER TABLE customers DROP COLUMN IF EXISTS devicelimit;
ALTER TABLE customers DROP COLUMN IF EXISTS expirytime;
ALTER TABLE customers DROP COLUMN IF EXISTS registrationtime;
ALTER TABLE customers DROP COLUMN IF EXISTS customerstatus;
ALTER TABLE customers DROP COLUMN IF EXISTS accounttype;
ALTER TABLE customers DROP COLUMN IF EXISTS email;

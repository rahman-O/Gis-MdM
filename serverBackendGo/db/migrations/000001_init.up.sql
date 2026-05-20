-- Minimal Headwind MDM schema for serverBackendGo auth (dev seed).
-- Column names match Liquibase; PostgreSQL folds unquoted identifiers to lowercase.

CREATE TABLE IF NOT EXISTS customers (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    master      BOOLEAN NOT NULL DEFAULT FALSE,
    filesdir    VARCHAR(200),
    prefix      VARCHAR(20) NOT NULL DEFAULT 'hmdm-',
    lastlogintime BIGINT
);

CREATE TABLE IF NOT EXISTS permissions (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    superadmin  BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS userroles (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(50) NOT NULL,
    description TEXT,
    superadmin  BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS userrolepermissions (
    roleid        INT NOT NULL REFERENCES userroles (id) ON DELETE CASCADE,
    permissionid  INT NOT NULL REFERENCES permissions (id) ON DELETE CASCADE,
    PRIMARY KEY (roleid, permissionid)
);

CREATE TABLE IF NOT EXISTS users (
    id                    SERIAL PRIMARY KEY,
    login                 VARCHAR(30) NOT NULL UNIQUE,
    email                 VARCHAR(50),
    name                  VARCHAR(50),
    password              VARCHAR(64) NOT NULL,
    customerid            INT REFERENCES customers (id) ON DELETE CASCADE,
    userroleid            INT REFERENCES userroles (id) ON DELETE RESTRICT,
    authtoken             VARCHAR(40),
    passwordresettoken    VARCHAR(40),
    lastloginfail         BIGINT NOT NULL DEFAULT 0,
    alldevicesavailable   BOOLEAN NOT NULL DEFAULT TRUE,
    allconfigavailable    BOOLEAN NOT NULL DEFAULT TRUE,
    passwordreset         BOOLEAN NOT NULL DEFAULT FALSE,
    twofactorsecret       VARCHAR(100),
    twofactoraccepted     BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS groups (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS configurations (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS userdevicegroupsaccess (
    userid  INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    groupid INT NOT NULL REFERENCES groups (id) ON DELETE CASCADE,
    PRIMARY KEY (userid, groupid)
);

CREATE TABLE IF NOT EXISTS userconfigurationaccess (
    userid            INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    configurationid   INT NOT NULL REFERENCES configurations (id) ON DELETE CASCADE,
    PRIMARY KEY (userid, configurationid)
);

CREATE TABLE IF NOT EXISTS settings (
    id           SERIAL PRIMARY KEY,
    customerid   INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    twofactor    BOOLEAN NOT NULL DEFAULT FALSE,
    idlelogout   INT
);

INSERT INTO customers (id, name, description, master, prefix)
VALUES (1, 'ADMIN', 'Master customer', TRUE, 'hmdm-')
ON CONFLICT (id) DO NOTHING;

INSERT INTO permissions (id, name, description, superadmin)
VALUES (1, 'superadmin', 'Super admin functions', TRUE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO userroles (id, name, description, superadmin)
VALUES
    (1, 'Super Admin', 'Full access', TRUE),
    (2, 'Organization Admin', 'Customer admin', FALSE),
    (3, 'User', 'Standard user', FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO userrolepermissions (roleid, permissionid)
VALUES (1, 1)
ON CONFLICT DO NOTHING;

INSERT INTO settings (customerid, twofactor, idlelogout)
SELECT 1, FALSE, NULL
WHERE NOT EXISTS (SELECT 1 FROM settings WHERE customerid = 1);

-- password: MD5("admin") uppercase hex -> SHA1(md5 + "5YdSYHyg2U") per Java PasswordUtil
INSERT INTO users (login, email, name, password, customerid, userroleid)
VALUES (
    'admin',
    'admin@localhost',
    'Administrator',
    '349242D38ED8667B5C11D2412EBEA4636BD3CA3A',
    1,
    1
)
ON CONFLICT (login) DO NOTHING;

SELECT setval('customers_id_seq', GREATEST((SELECT MAX(id) FROM customers), 1));
SELECT setval('userroles_id_seq', GREATEST((SELECT MAX(id) FROM userroles), 3));
SELECT setval('users_id_seq', GREATEST((SELECT MAX(id) FROM users), 1));

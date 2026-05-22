CREATE TABLE IF NOT EXISTS usagestats (
    id SERIAL PRIMARY KEY,
    ts DATE NOT NULL DEFAULT CURRENT_DATE,
    instanceid VARCHAR(255),
    webversion VARCHAR(255),
    community BOOLEAN NOT NULL DEFAULT TRUE,
    devicestotal INT NOT NULL DEFAULT 0,
    devicesonline INT NOT NULL DEFAULT 0,
    cputotal INT NOT NULL DEFAULT 0,
    cpuused INT NOT NULL DEFAULT 0,
    ramtotal INT NOT NULL DEFAULT 0,
    ramused INT NOT NULL DEFAULT 0,
    scheme VARCHAR(255),
    arch VARCHAR(255),
    os VARCHAR(255),
    UNIQUE (ts, instanceid)
);

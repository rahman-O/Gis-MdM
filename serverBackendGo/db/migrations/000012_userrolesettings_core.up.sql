CREATE TABLE IF NOT EXISTS userrolesettings (
    id SERIAL PRIMARY KEY,
    roleid INT NOT NULL REFERENCES userroles (id) ON DELETE CASCADE,
    customerid INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    columndisplayeddevicestatus BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayeddevicedate BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayeddevicenumber BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayeddevicemodel BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayeddevicepermissionsstatus BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayeddeviceappinstallstatus BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayeddeviceconfiguration BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayeddeviceimei BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayeddevicephone BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayeddevicedesc BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayeddevicegroup BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayedlauncherversion BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayeddevicefilesstatus BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayedbatterylevel BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayeddefaultlauncher BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayedcustom1 BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayedcustom2 BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayedcustom3 BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayedmdmmode BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayedkioskmode BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayedandroidversion BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayedenrollmentdate BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayedserial BOOLEAN NOT NULL DEFAULT TRUE,
    columndisplayedpublicip BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (roleid, customerid)
);

INSERT INTO userrolesettings (roleid, customerid)
SELECT ur.id, c.id
FROM userroles ur
CROSS JOIN customers c
WHERE ur.superadmin = FALSE
  AND NOT EXISTS (
      SELECT 1 FROM userrolesettings urs
      WHERE urs.roleid = ur.id AND urs.customerid = c.id
  );

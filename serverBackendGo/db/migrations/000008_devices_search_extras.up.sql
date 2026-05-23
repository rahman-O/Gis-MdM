-- 012 US1: optional columns for Java-aligned device search filters.

ALTER TABLE devices ADD COLUMN IF NOT EXISTS imeiupdatets BIGINT;

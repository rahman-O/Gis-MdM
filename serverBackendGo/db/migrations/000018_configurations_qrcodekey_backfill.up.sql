-- Backfill missing enrollment QR keys (Java used NOT NULL DEFAULT MD5(RANDOM()::TEXT)).
UPDATE configurations
SET qrcodekey = md5(random()::text)
WHERE qrcodekey IS NULL OR trim(qrcodekey) = '';

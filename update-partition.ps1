param(
  [int]$OldId = 1,
  [string]$NewTitle = "Partition PRO (double gypsum + wool)",
  [string]$Category = "Professional",
  [string]$Description = "Double gypsum 12.5mm + mineral wool 50mm + frame 50mm",
  [string]$NoiseReduction = "52-56 dB",
  [string]$Thickness = "12-14 cm",
  [string]$Material = "Gypsum + mineral wool",
  [string]$PricePerSqm = "1200-1600 RUB/m2",
  [string]$ImageUrl = ""
)

$ErrorActionPreference = "Stop"

$updateSql = @"
BEGIN;
UPDATE partitions SET is_active = FALSE WHERE id = $OldId;
INSERT INTO partitions (
  title, category, description, noise_reduction, thickness, material, price_per_sqm, image_url, is_active
) VALUES (
  '$NewTitle', '$Category', '$Description', '$NoiseReduction', '$Thickness', '$Material', '$PricePerSqm', NULLIF('$ImageUrl',''), TRUE
);
COMMIT;
"@

$updateSql | docker exec -i partition_postgres psql -U partition_user -d partition -v ON_ERROR_STOP=1

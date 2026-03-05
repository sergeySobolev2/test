param(
  [int]$Id = 1
)

$ErrorActionPreference = "Stop"

$sql = @"
BEGIN;
DELETE FROM calculation_items WHERE partition_id = $Id;
DELETE FROM partitions WHERE id = $Id;
COMMIT;
"@

$sql | docker exec -i partition_postgres psql -U partition_user -d partition -v ON_ERROR_STOP=1

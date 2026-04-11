-- Backfill tblUserBranches_UBR from tblUserRoles_URO.
-- Users with role assignments on branches were never added to the
-- visibility table, so apps that query tblUserBranches_UBR (VetData,
-- OftaData) couldn't resolve their branches.
INSERT INTO "tblUserBranches_UBR" (usr_id, bra_id, created_at)
SELECT DISTINCT ur.usr_id, ur.bra_id, NOW()
FROM "tblUserRoles_URO" ur
WHERE NOT EXISTS (
    SELECT 1 FROM "tblUserBranches_UBR" ub
    WHERE ub.usr_id = ur.usr_id AND ub.bra_id = ur.bra_id
);

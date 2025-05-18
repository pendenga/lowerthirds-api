SELECT *
FROM Users u
INNER JOIN OrgUsers ou
  ON ou.user_id = u.id
INNER JOIN Organization o
  ON o.id = ou.org_id
WHERE u.email = 'pendenga@gmail.com';


-- user map for active users and where the calling user has access to those orgs
SELECT ou.org_id, ou.user_id
FROM OrgUsers ou
INNER JOIN Users u
  ON u.id = ou.user_id
  AND u.deleted_dt IS NULL
INNER JOIN Organization o
  ON o.id = ou.org_id
  AND o.deleted_dt IS NULL
WHERE ou.deleted_dt IS NULL
  AND ou.org_id IN (SELECT org_id FROM OrgUsers WHERE user_id = '3cd5fe4e-9ecb-4ec2-b7c7-0d19288c08e0')

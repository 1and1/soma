package stmt

const SelectUserToken = `
SELECT salt,
       valid_from,
       valid_until
FROM   auth.user_token_authentication
WHERE  token = $1::varchar;`

const SelectAdminToken = `
SELECT salt,
       valid_from,
       valid_until
FROM   auth.user_token_authentication
WHERE  token = $1::varchar;`

const SelectToolToken = `
SELECT salt,
       valid_from,
       valid_until
FROM   auth.user_token_authentication
WHERE  token = $1::varchar;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix

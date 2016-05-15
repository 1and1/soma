package stmt

const SelectToken = `
SELECT salt,
       valid_from,
       valid_until
FROM   auth.tokens
WHERE  token = $1::varchar;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix

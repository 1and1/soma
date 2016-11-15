/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const (
	SupervisorCategoryStatements = ``

	CategoryAdd = `
INSERT INTO soma.categories (
            category,
            created_by
)
SELECT $1::varchar,
       ( SELECT user_id
         FROM   inventory.users
         WHERE  user_uid = $2::varchar)
WHERE NOT EXISTS (
      SELECT category
      FROM   soma.categories
      WHERE  category = $1::varchar);`

	CategoryRemove = `
DELETE FROM soma.categories
WHERE category = $1::varchar;`

	CategoryList = `
SELECT category
FROM   soma.categories;`

	CategoryShow = `
SELECT sc.category,
       iu.user_uid,
       sc.created_at
FROM   soma.categories sc
JOIN   inventory.users iu
ON     sc.created_by = iu.user_id
WHERE  sc.category = $1::varchar;`
)

func init() {
	m[CategoryAdd] = `CategoryAdd`
	m[CategoryList] = `CategoryList`
	m[CategoryRemove] = `CategoryRemove`
	m[CategoryShow] = `CategoryShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix

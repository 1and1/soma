/*-
 * Copyright (c) 2016, Jörg Pernfuß <code.jpe@gmail.com>
 * No rights reserved.
 *
 * This code is available licensed as CC0. Please see the included
 * LICENSE file for more information.
 */

package hash64

import "encoding/base64"

const encodeHash64 string = `./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`

var StdEncoding = base64.NewEncoding(encodeHash64).WithPadding(base64.NoPadding)
var PadEncoding = base64.NewEncoding(encodeHash64)

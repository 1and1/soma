Modular crypt format specification for type scrypt-h64
======================================================

The following document specifies the format of the modular crypt
format (MCF) for identifier `scrypt-h64`.

Key derivation function
=======================

The used key derivation function is `scrypt` utilizing
`PBKDF2-HMAC-SHA256`.

Field Separator
===============

The used field seprator must be the character `$`.

The MCF string must start with a leading field separator.

The MCF string must not end with a trailing field separator.

Implementations must handle a trailing end field separator by
stripping it.

Fields
======

The following fields are used, in the specified order:

1. Identifier
2. Parameters
3. Salt
4. Digest

If a field is empty, its field separator must not be skipped.

Field: Identifier
=================

The value of the field identifier is `scrypt-h64`.

This field must not be empty.

Field: Parameters
=================

This field may contain the following parameters:

1. N
2. r
3. p
4. l
5. s

Parameters are encoded `key=value` with no whitespace.

Multiple parameters are separated by a `,` with no surrounding
whitespace.

Parameters may be ordered in any order.

Any paramter may be omitted, up to an empty parameter field.

Parameter: N
------------

The value of parameter `N` is the log2 of scrypts work factor N.

Valid value range is 1-65535.

If `N` is not specified, the default value is `14`.

Parameter: r
------------

The value of parameter `r` is scrypts block size factor r.

Valid value range is 1-255.

If `r` is not specified, the default value is `8`.

Parameter: p
------------

The value of parameter `p` is scrypts parallel work factor p.

Valid value range is 1-255.

If `p` is not specified, the default value is `1`.

Parameter: l
------------

The value of paramater `l` is scrypts digest output length dkLen,
setting the length of the derived key in bytes.

Valid value range is 16-65535.

If `l` is not specified, the default value is `32` for a 256 bit
digest.

Parameter: s
------------

The value of parameter `s` is the length of the used salt in bytes.

This parameter is specified, so no implicit assumptions about the
length of the used salt are made.

Valid value range is 16-65535.

If `s` is not specified, the default value is `16` for a 128 bit
salt.

Field: Salt
===========

This field contains the Hash64 encoded representation of the salt.

Field: Digest
=============

This field contains the Hash64 encoded representation of the digest.

Appendix: Hash64
================

Hash64 is a base64 encoding variation using the following character
sequence:

```
./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz
```
This scheme does not use padding.

The name has been coined by passlib. It is the traditional encoding
map used for `des_crypt`, `md5_crypt $1$` and others.

Appendix: Examples
==================

Used password: `correct horse battery staple`

Bad example
-----------

This example uses `4MiByte` RAM.

```
Parameter{
    Iterations: 12,
    Length:     16,
    SaltLength: 16,
}

% time ./scrypth64
$scrypt-h64$N=12,r=8,p=1,l=16,s=16$t3QnR5Ck2KVlkkK5zqjZZU$m.a/EOXM/RbQ3q9ghFqEI.
./scrypth64  0.02s user 0.00s system 101% cpu 0.023 total
```

Good example
------------

This example uses `64MiByte` RAM.

```
DefaultParams()

% time ./scrypth64
$scrypt-h64$N=15,r=16,p=2,l=48,s=64$gSBRS/x9K5aguQLY4X90/P6hPMoC20K2LOSYajzDObyIzeg3K4YxMyOlA3/FGSK1LBKD2hTxrWI2UbBDHhD3pE$SY7Qed/M.1SnnQL8aeO6850MV5bQSWpxzmThhmOz7eu0MkK/EM4rdaS4C0Yt1iOj
./scrypth64  0.58s user 0.02s system 100% cpu 0.604 total
```

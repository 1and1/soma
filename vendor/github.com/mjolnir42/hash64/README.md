Hash64 - Encoding Wrapper around base64
=======================================

`Hash64` is the term Python's `passlib` coined for a variant of base64
encoding traditionally used for encoding password hashes in
`passwd(5)` files.

FreeBSD's md5crypt by Poul-Henning Kamp introduced the `$1$`
identifier which was chosen since `$` does not appear in the Hash64
character map and thus avoided collisions with existing DES hashes;
thus laying the groundwork for what is now refered to as the Modular
Crypt Format. A number of schemes in MCF use Hash64 to encode salt and
hash digest.

This encoding type is usually unpadded.

Hash64 implements two trivial wrappers `Hash64.StdEncoding` for the
unpadded encoding and `Hash64.PadEncoding` for the padded one. They
are of type `*base64.Encoding`, and as such the regular functions
offered by it can be used: `Decode`, `DecodeToString`, `DecodedLen`,
`Encode`, `EncodeToString` and `EncodedLen`.

License
-------

It feels silly to add a license to 3 lines of code, but no license
means no usage at all. This code is therefor put under Creative
Commons Zero.

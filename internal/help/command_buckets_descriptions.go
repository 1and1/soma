package help

var CmdBucketPropertyAdd string = `
   Add a property to a bucket.

SYNOPSIS:
   somaadm buckets property add \
      ${propertyName} \
      type ${propertyType} \
      to ${bucketName} \
      [ value ${propertyValue} ]

ARGUMENT TYPES:

EXAMPLES:
   somaadm buckets property add \
      'DNS Resolver [unbound]' \
      type service
      to itolive`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix

# somaadm buckets property add

Add a property to a bucket

# SYNOPSIS
```
somaadm buckets property add \ 
    ${propertyName} \
    type ${propertyType} \
    to ${bucketName} \
    [ value ${propertyValue} ]
```

# ARGUMENT TYPES

# EXAMPLES
```
somaadm buckets property add \
    'DNS Resolver [unbound]' \
    type service
    to common_master
```

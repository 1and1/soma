# somaadm servers create

This command creates a new physical server in SOMA, allowing
to set all server attributes except the deleted flag; a server
can not be created as deleted.
The online argument is optional and defaults to true. The first
argument is the name of the server. The order of additional
argument/value pairs is irrelevant.

# SYNOPSIS

```
somaadm servers create ${name} \
   assetid ${assetId} \
   datacenter ${datacenter} \
   location ${location} \
   [ online ${online} ]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | -------- 
name | string | Name of the server ("hostname") | | no
assetid | uint64 | Numeric id of the server in ext. asset tracking | | no
datacenter | string | Datacenter the server is located in | | no
location | string | Location of the server within the datacenter | | no
online | bool | Server is online or not | true | yes

# PERMISSIONS

# EXAMPLES

```
./somaadm servers create 'eu-west-1a' datacenter 'aws.eu-west-1' \
   location 'Zone A' assetid 12345 online false

./somaadm servers create 'eu-central-1b' assetid 6349 \
   datacenter 'aws.eu-central-1' location 'Zone B'

./somaadm servers create 'fooserver' assetid 42 \
   datacenter 'de.fra' location 'Room 3 Rack 91 Slot 4'
```

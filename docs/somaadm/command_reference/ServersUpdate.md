# somaadm servers update

This command updates all attributes of a physical server. Exept for
the id, all values are replaced.

# SYNOPSIS

```
somaadm servers update ${uuid} \
   name ${newname} \
   assetid ${assetid} \
   datacenter ${datacenter} \
   location ${location} \
   online ${online} \
   deleted ${deleted}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | -------- 
uuid | string | UUID of the server to update | | no
newname | string | New name of the server ("hostname") | | no
assetid | uint64 | Numeric id of the server in ext. asset tracking | | no
datacenter | string | Datacenter the server is located in | | no
location | string | Location of the server within the datacenter | | no
online | bool | Server is online or not | | no
deleted | bool | Server is marked deleted or not | | no

# PERMISSIONS

# EXAMPLES

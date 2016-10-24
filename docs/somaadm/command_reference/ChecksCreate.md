# somaadm checks create

This command creates a new check configuration in SOMA. Every check configuration
is named, and the name is unique per repository.
Checks can be created on repositories, buckets, groups, clusters and nodes.
However only groups, clusters and nodes will compute check instances for
deployment, ie. you can not alert your configuration repository metadata.

Checks are defined on buckets, since groups and clusters are unique
per bucket. Checks on repositories still have a bucket as argument to 'in'.

A check must have at least one threshold but can have as many thresholds as
the capability allows. There are no set limitations to the amount of
constraints a check may have.

This command is asynchronous and returns a JobID.

# SYNOPSIS

```
somaadm checks create ${check} \
   in ${bucket} \
   on ${type} ${object} \
   with ${capability} \
   interval ${intv} \
   threshold predicate ${symbol} level ${lvl} value ${val} \
     [ [ threshold ... ] ... ] \
   [ inheritance ${inherit} ] \
   [ childrenonly ${child} ] \
   [ extern ${extid} ] \
   [ [ constraint ${ctype} ${prop} ${cval} ] \
     [ constraint ... ] ... ]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | -------- 
check | string | Name of the check configuration | | no
bucket | string | Name of a bucket in the repository | | no
type | string | Type of the object the check is on | | no
object | string | Name of the object to receive the check | | no
capability | string | The capability to use for this check | | no
intv | uint64 | Checkinterval in seconds, greater 0 | | no
inherit | bool | Check is inherited to child objects | true | yes
child | bool | Check is only active on child objects | false | yes
extid | string | External correlation id | | yes
symbol | string | Predicate symbol to compare the threshold with | | no
lvl | string | Name of the notification level to alert at | | no
val | string | Threshold value | | no
ctype | string | Property type to constraint against | | no
prop | string | The property to constraint against | | no
cval | string | Value to constraint against. '@defined' acts as magic accepting all values. | | no

# PERMISSIONS

# EXAMPLES

```
./somaadm checks create 'default node ping'      \
   in common_master                              \
   on repository common                          \
   with icinga.internal.icmp.rtt                 \
   threshold predicate '>=' level info value 450 \
   interval 300                                  \
   constraint native object_type node
```

This creates a check configuration on the repository that is fully inherited
but only evaluated on leaf objects of type node. The check runs every
5 minutes and alerts at informational level if the measured RTT is over
450 milliseconds, the unit of metric icmp.rtt.

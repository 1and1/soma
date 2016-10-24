# somaadm checks show

This command shows a specific check configuration. Check names are unique per
repository. They are adressed via a bucket to keep the cli syntax consistent
with regards to the `in` keyword.

# SYNOPSIS

```
somaadm checks show ${check} in ${bucket}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | -------- 
check | string | Name of the check configuration | | no
bucket | string | Name of a bucket in the repository | | no

# PERMISSIONS

# EXAMPLES

```
./somaadm checks show 'default node ping' in common_master
```

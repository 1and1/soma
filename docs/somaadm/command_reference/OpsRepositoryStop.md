# somaadm ops repository stop

This command stops a specific repository on the server. The
repository has to be either in state 'ready' or 'broken',
ie. fully operational or faulted.

A running repository that is still loading can not be stopped.

# SYNOPSIS

```
somaadm ops repository stop ${repository}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | -------- 
repository | string | Name or UUID of the repository to stop | | no

# PERMISSIONS

This command requires one of the following permissions:

* system\_all

# EXAMPLES

```
./somaadm ops repository stop common

./somaadm ops repository stop 21a4eda3-fafd-4c1b-9dd2-a29ef96ac916
```

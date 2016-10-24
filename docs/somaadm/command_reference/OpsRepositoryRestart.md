# somaadm ops repository restart

This command restarts the running or stopped repository on the
server.

A running repository that is still loading can not be stopped.
A broken repository can be restarted in the hope that the
problem goes away (it won't).

# SYNOPSIS

```
somaadm ops repository restart ${repository}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | -------- 
repository | string | Name or UUID of the repository to restart | | no

# PERMISSIONS

This command requires one of the following permissions:

* system\_all

# EXAMPLES

```
./somaadm ops repository restart common

./somaadm ops repository restart 21a4eda3-fafd-4c1b-9dd2-a29ef96ac916
```

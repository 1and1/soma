# somaadm checks delete

This command deletes a check definition from a repository. Check names are
unique per repository. They are addressed via a bucket to keep the cli syntax
consistent with regards to the 'in' keyword.

Deleting a check definition also deletes all checks derived from the definition,
direct and inherited, as well as all check instances these checks spawned.
All currently deployed check instances will be deprovisioned.

This command is asynchronous and returns a JobID.

# SYNOPSIS

```
somaadm checks delete ${check} in ${bucket}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | -------- 
check | string | Name of the check definition | | no
bucket | string | Name of a bucket in the repository | | no

# PERMISSIONS

# EXAMPLES

```
./somaadm checks delete 'default node ping' in common_master
```

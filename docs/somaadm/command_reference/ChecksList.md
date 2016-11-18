# somaadm checks list

This command lists all check definitions in a repository. Check names are
unique per repository. They are addressed via a bucket to keep the cli syntax
consistent with regards to the `in` keyword.

# SYNOPSIS

```
somaadm checks list in ${bucket}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | -------- 
bucket | string | Name of a bucket in the repository | | no

# PERMISSIONS

# EXAMPLES

```
./somaadm checks list in common_master
```

This will list all check configurations in repository common. It is
currently not possible to list only the check configurations of a
subtree.

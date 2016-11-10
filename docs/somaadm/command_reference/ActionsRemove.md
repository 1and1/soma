# somaadm actions remove

This command is used to delete an action from a section.

# SYNOPSIS

```
somaadm actions remove ${action} from ${section}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
action | string | Name of the action | | no
section | string | Name of the section | | no

# PERMISSIONS

# EXAMPLES

```
./somaadm actions remove add from environments
./somaadm actions remove remove from environments
./somaadm actions remove list from environments
./somaadm actions remove show from environments
```

# somaadm actions add

This command is used to add a new permission action to a permission
section.

Actions are grouped in sections, and are used to build permissions
from.

# SYNOPSIS

```
somaadm actions add ${action} to ${section}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
action | string | Name of the action | | no
section | string | Name of the section | | no

# PERMISSIONS

# EXAMPLES

```
./somaadm actions add add to environments
./somaadm actions add remove to environments
./somaadm actions add list to environments
./somaadm actions add show to environments
```

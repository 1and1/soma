# somaadm capabilities declare

This command creates a new capability, which is a declaration
what a monitoring system can monitor. It is defined as the
combination of which monitoring system provides the service, which
view is utilized, which metric is monitored and how many threshold
definitions can be used.

Optionally, demux specifications can be provided; referencing a service
attribute. If a servie has this attribute more than once, for every
instance of it a separate check instance will be created.

Optionally, constraints can be specified. These constraints must be
fulfilled for the resulting check to work correctly, ie. for a check
instance to be created, these have to be fulfilled as well.
As a special constraint value, the string `@defined` may be used to
indicate that the value is not important, but it has to be set.

# SYNOPSIS

```
somaadm capabilities \
   declare ${monitoring} \
   view ${view} \
   metric ${metric} \
   thresholds ${count} \
   [ [ demux ${attribute} ],
     ... ] \
   [ [ constraint ${type} ${name} ${value} ],
     ... ]

```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | -------- 
monitoring | string | Monitoring system to use | | no
view | string | View to use for monitoring | | no
metric | string | Monitored metric path specification | | no
count | uint64 | More thresholds than levels is invalid | | no
attribute | string | Valid service attribute name | |
type | string | Property type to constraint against | |
name | string | Name of the property | | 
value | string | Value of the property | |

# PERMISSIONS

# EXAMPLES

```
./somaadm capabilities declare icinga \
   view internal \
   metric tcp.rtt \
   thresholds 1 \
   constraint attribute port @defined \
   constraint attribute transport_protocol tcp \
   demux port`
```

# BUGS

Capability constraints and demux definitions are currently
not implemented.

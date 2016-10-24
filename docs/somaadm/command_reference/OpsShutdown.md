# somaadm ops shutdown

This command will activate the controlled shutdown of a SOMA
instance. Once received, it will answer all requests with a HTTP
unavailable error, while finishing active jobs, draining queues
and shutting down background go routines.

This command has no confirmation dialogue.

# SYNOPSIS

```
somaadm ops shutdown
```

# ARGUMENT TYPES

# PERMISSIONS

This command requires one of the following permissions:

* system\_all

# EXAMPLES

```
./somaadm ops shutdown
```

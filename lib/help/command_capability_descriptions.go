package help

var CmdCapabilityDeclare string = `
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
	 As a special constraint value, the string "@defined" may be used to
	 indicate that the value is not important, but it has to be set.

SYNOPSIS:
   somaadm capabilities \
	    declare ${monitoring} \
      view ${view} \
      metric ${metric} \
      thresholds ${count} \
			[ [ demux ${attribute} ],
			  ... ] \
			[ [ constraint system|attribute ${system|attribute} ${value} ],
			  ... ]

ARGUMENT TYPES:
   monitoring string    Monitoring system to use
	 view       string    View to use for monitoring
	 metric     string    Monitored metric path specification
	 count      uint64    More thresholds than levels is invalid
	 attribute  string    Valid service attribute name
	 system     string    Valid system property name

EXAMPLES:
   ./somaadm capabilities declare icinga \
	    view internal \
	    metric tcp.rtt \
			thresholds 1 \
			constraint attribute port @defined \
			constraint attribute transport_protocol tcp \
			demux port`

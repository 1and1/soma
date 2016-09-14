package help

const CmdCheckAdd string = `
   This command creates a new check definition in SOMA. Every check definition
   is named, and the name is unique per repository.
   Checks can be created on repositories, buckets, groups, clusters and nodes.
   However only groups, clusters and nodes will compute check instances for
   deployment, ie. you can not alert your configuration repository metadata.

   Checks are defined on buckets, since groups and clusters are unique
   per bucket. Checks on repositories still have a bucket as argument to 'in'.

   A check can have as many thresholds as the capability allows.
   There no practical limitations to the amount of constraints a check may have.

   This command is asynchronous and returns a JobID.

SYNOPSIS:
   somaadm checks create ${check} \
      in ${bucket} \
      on ${type} ${object} \
      with ${capability} \
      interval ${intv} \
      [ inheritance ${inherit} ] \
      [ childrenonly ${child} ] \
      [ extern ${extid} ] \
      [ [ threshold predicate ${symbol} level ${lvl} value ${val} ]
        , ... ] \
      [ [ constraint ${ctype} ${prop} ${cval} ]
        , ... ]

ARGUMENT TYPES:
   check      string    Name of the check definition
   bucket     string    Name of a bucket in the repository
   type       string    Type of the object the check is on, ie. one of
                        repository, bucket, group, cluster or node
   object     string    Name of the object 
   capability string    The capability to use for this check
   intv       uint64    Checkinterval in seconds, > 0
   inherit    bool      Check is inherited to child objects. Default: true
   child      bool      Check is only active on child objects. Default: false
   extid      string    External correlation id that can be set
   symbol     string    Predicate symbol to compare the threshold with
   lvl        string    Name of the notification level to alert at
   val        string    Threshold value
   ctype      string    Property type to constraint against. Possible values:
                        service, oncall, attribute, system, native, custom
   prop       string    The property to constraint against. For type service
                        and oncall, only 'name' is available
   cval       string    Value to constraint against. '@defined' acts as magic
                        accepting all values.

EXAMPLE:
   ./somaadm checks create 'default broadcast ping' \
      in common_master                              \
      on repository common                          \
      with icinga.internal.tcp.rtt                  \
      threshold predicate '>=' level info value 450 \
      interval 300                                  \
      constraint native object_type node

   This creates a check definition on the repository that is fully inherited
   but only evaluated on leaf objects of type node. The check runs every
   5 minutes and alerts at informational level if the measured RTT is over
   450 milliseconds, the unit of metric tcp.rtt.
`

const CmdCheckDelete string = `
   This command deletes a check definition from a repository. Check names are
   unique per repository. They are adressed via a bucket to keep the cli syntax
   consistent with regards to the 'in' keyword.

   Deleting a check definition also deletes all checks derived from the definition,
   direct and inherited, as well as all check instances these checks spawned.
   All currently deployed check instances will be deprovisioned.

   This command is asynchronous and returns a JobID.

SYNOPSIS:
   somaadm checks delete ${check} in ${bucket}

ARGUMENT TYPES:
   check      string    Name of the check definition
   bucket     string    Name of a bucket in the repository

EXAMPLES:
   ./somaadm checks delete 'Default ping' in foobar_default
`

const CmdCheckList string = `
   This command lists all check definitions in a repository. Check names are
   unique per repository. They are adressed via a bucket to keep the cli syntax
   consistent with regards to the 'in' keyword.

SYNOPSIS:
   somaadm checks list in ${bucket}

ARGUMENT TYPES:
   bucket     string    Name of a bucket in the repository

EXAMPLES:
   ./somaadm checks list in foobar_default
`

const CmdCheckShow string = `
   This command shows a specific check. Check names are unique per repository.
   They are adressed via a bucket to keep the cli syntax consistent with regards
   to the 'in' keyword.

SYNOPSIS:
   somaadm checks show ${check} in ${bucket}

ARGUMENT TYPES:
   check      string    Name of the check definition
   bucket     string    Name of a bucket in the repository

EXAMPLES:
   ./somaadm checks show 'Default ping' in foobar_default
`

package help

var CmdServerCreate string = `
   This command creates a new physical server in SOMA, allowing
   to set all server attributes except the deleted flag; a server
   can not be created as deleted.
   The online argument is optional and defaults to true. The first
   argument is the name of the server. The order of additional
   argument/value pairs is irrelevant.

SYNOPSIS:
   somaadm servers create ${name} \
      assetid ${assetId} \
      datacenter ${datacenter} \
      location ${location} \
      [ online ${online} ]

ARGUMENT TYPES:
   name       string    Name of the server ("hostname")
   assetid    uint64    Numeric id of the server in ext. asset tracking
   datacenter string    Datacenter the server is located in
   location   string    Location of the server within the datacenter
   online     bool      Server is online or not

EXAMPLES:
   ./somaadm servers create 'fooserver' datacenter 'de.ber' \
      location 'Slot 69' id 12345

   ./somaadm servers create 'barserver' id 6349 \
      datacenter 'de.fra.fra' location 'Room 3 Unit 12' online false`

var CmdServerUpdate string = `
   This command updates all attributes of a physical server. Exept for
   the id, all values are replaced.

SYNOPSIS:
   somaadm servers update ${uuid} \
      name ${newname} \
      assetid ${assetid} \
      datacenter ${datacenter} \
      location ${location} \
      online ${online} \
      deleted ${deleted}

ARGUMENT TYPES:
   uuid       string    UUID of the server to update
   newname    string    New name of the server ("hostname")
   assetid    uint64    Numeric id of the server in ext. asset tracking
   datacenter string    Datacenter the server is located in
   location   string    Location of the server within the datacenter
   online     bool      Server is online or not
   deleted    bool      Server is marked deleted or not`

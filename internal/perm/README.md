# perm

```
package perm // import "github.com/1and1/soma/internal/perm"

Package perm implements the permission cache module for the SOMA supervisor.
It tracks which actions are mapped to permissions and which permissions have
been granted.

It can be queried whether a given user is authorized to perform an action.

type Cache struct { ... }
    func New() *Cache
```

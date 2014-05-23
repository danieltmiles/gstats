# Package log

```Go
import "github.com/monsooncommerce/log"
```

This package formats logs in the following format:

```
TIMESTAMP HOSTNAME TAG[PID]: LEVEL MESSAGE
```

Timestamp is RFC3339 formated. Hostname is the FQDN of the host writing the message.
Tag represents the application name. PID is the Unix PID of the running application.
Level is one of DEBUG, ERROR, INFO, NOTICE, or WARNING. All messages are written to
stdout except ERROR messages, which are written to stderr.

## Usage

```Go
package main

import (
    "github.com/monsooncommerce/log"
)

func main() {
    log.SetTag("myapp")
    log.Error("This is an error message")
}
```

Output:

```
2013-10-02T13:41:43Z host.monsooncommerce.com myapp[25702]: ERROR This is an error message
```

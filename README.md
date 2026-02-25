# portmap

Port Scanner - lists all listening ports with process information.

## Purpose

Display all TCP and UDP ports in use on the system along with associated process details.

## Installation

```bash
go build -o portmap ./cmd/portmap
```

## Usage

```bash
portmap
```

### Examples

```bash
# List all listening ports
portmap

# Combined with other tools
portmap | grep 8080  # Find port 8080
```

## Output

```
=== PORT MAP ===

8080   tcp    0.0.0.0:8080      node
    PID: 12345 | node server.js

3306   tcp    127.0.0.1:3306      mysqld
    PID: 23456 | /usr/sbin/mysqld

22     tcp    0.0.0.0:22        sshd
    PID: 1111 | /usr/sbin/sshd -D
```

## Port Color Coding

- Red: Well-known ports (<1024)
- Yellow: Registered ports (1024-7999)
- Green: Dynamic/private ports (8000+)

## Dependencies

- Go 1.21+
- github.com/fatih/color
- ss (from iproute2 package)
- lsof
- ps

## Build and Run

```bash
# Build
go build -o portmap ./cmd/portmap

# Run
go run ./cmd/portmap
```

## Useful Commands

```bash
lsof -i :<port>     # Find process using specific port
kill <pid>          # Kill process by PID
ss -tuln            # List all listening ports
```

## License

MIT
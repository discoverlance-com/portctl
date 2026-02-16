# Ports Manager (portctl)

`portctl` is a lightweight cross platform CLI for identifying and terminating processes listening on network ports without having to remember OS-specific commands like `netstat`, `lsof`, or `taskkill`.

## Key Features

- **True Cross-Platform Support**: Cross-platform support for Windows, Linux and macOS using system native tools.
- **Safety First**: Built-in confirmation prompt when terminating processes to avoid accidental process termination.
- **Flexible**: Kill process by either port or process id.
- **Zero Dependencies**: Relies strictly on the Go standard library and native OS tools - no heavy external framework required.

## Getting Started

### Installation & Setup

#### Prerequisite

- Go 1.25.1 or higher installed on your system.
- Terminal (Powershell, Bash, Zsh)

#### Install Via Go

Install the CLI directly with Go.

```bash
go install github.com/discoverlance-com/portctl@latest
```

#### Build From Source

| Platform | Command |
| :--- | :--- |
| Windows | `$env:GOOS="windows"; go build -o portctl.exe main.go` |
| Linux | `GOOS=linux go build -o portctl main.go` |
| macOS | `GOOS=darwin go build -o portctl main.go` |


#### Global Access (Optional)

To run the tool from anywhere, move the binary to a folder in your system's PATH.

- Linux/macOS: sudo mv portctl /usr/local/bin/
- Windows: Move portctl.exe to a folder like C:\bin and add that folder to your Environment Variables.


### Usage Guide

Let's see how we can interact with `portctl`.

#### List Listening Processes 

Find all services currently listening on network ports.

```bash
portctl list
```

#### Kill A Process 

Find all services currently listening on network ports.

```bash
# By Process ID
portctl kill -pid 1234

# By Port number
portctl kill -port 3000

# Force Kill (no prompt)
portctl kill -port 3000 -y
```


#### Example Output

```console
$ portctl list

PORT   PID
3000   4231
8080   4971
```

```console
$ portctl kill -help

Usage of kill:
  -port int
        The port of the running service you want to kill
  -pid int
        The process id for the running service you want to kill
  -y    Confirmation that you want to kill the process
```

```console
$ portctl kill -port 3000

Are you sure you want to continue? Y/N: y
Process: 4231 killed successfully
```

## Development

This project uses **Go build tags** and **Interfaces** to maintain a clean and testable codebase across different Operating systems.

- **Build Tags**: OS-specific logic is isolated in `windows.go`, `linux.go`, and `darwin.go` files.
- **The PortManager Interface**: A `PortManager` interface abstracts platform-specific implementations allowing the CLI commands to remain OS-agnostic.

## Contributing

Any **Contributions** are greatly **appreciated**. The goal of this project is to remain lightweight and depend primarily on the Go standard library. External dependencies should be justified by significant value or unavoidable complexity.

1. Fork the Project.
2. Create your Feature Branch (`git checkout -b feat/amazingupdate`).
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the Branch (`git push origin feat/amazingupdate`).
5. Open a Pull Request.

## License
Distributed under the MIT License. See [LICENSE](./LICENSE) for more information.

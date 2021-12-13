# WSL subnet utility

This is a small Go utility to set the WSL2 host and subnet. It achieves this by:
- deleting the existing WSL network
- creating a new one with the specified subnet (defaulting to `192.168.100.0/24`)

Windows automatically creates the `WSL` network when a WSL2 environment is started, so we need to pre-empt this by creating our own network with our settings before that happens.

## Important notes

- this needs to be done **before** starting any WSL2 (or Docker on WSL2) environments
- if you need to do this afterwards, restart WSL2 environments by shutting them down with `wsl --shutdown`; they should work correctly after they're re-launched.
- you may need to disable Docker auto-start: this needs to run first
- the utility needs to be run in an **elevated** mode (administrator) in order to alter the network configuration
  - you'll get an `Access denied` error if you are not elevated
  - run in an elevated console
  - or schedule as a task (noted below)

Basic command line help can be obtained with `wsl-subnet --help`

Installing this as a task is possible via the Task Scheduler:
- Use `At system startup` as the trigger
- Run under the `SYSTEM` account

This can done in an elevated console. Adjust the following example to point to your file location, of course.
```batch
schtasks /create /tn "WSL Subnet Configure" /tr c:\tools\wsl-subnet.exe /sc onstart /ru System
```

## Build

```batch
cd src
go build .
```

This will yield a single, static executable: `wsl-subnet.exe`


## References and acknowledgements

We're making use the Windows Host Compute Network interfaces to do this: 
- https://docs.microsoft.com/en-us/windows-server/networking/technologies/hcn/hcn-declaration-handles
- https://docs.microsoft.com/en-us/windows-server/networking/technologies/hcn/hcn-json-document-schemas

Fortunately, Microsoft has supplied the `hcsshim` Go library to interact with these interfaces at a higher level. This library is used by various projects, including (at some stage), Docker itself.
- https://pkg.go.dev/github.com/Microsoft/hcsshim

This is a simple utility inspired partly by the various Powershell scripts out there to control WSL booting and subnet assignment, including:
- https://github.com/ocroz/wsl2-boot
- https://github.com/skorhone/wsl2-custom-network
- https://github.com/jgregmac/hyperv-fix-for-devs
- https://github.com/wikiped/WSL-IpHandler

Check those out if you need something more complex. 

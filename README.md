# raft

Implementation of the Raft consensus protocol for distributed log consistency

## Use (Package)
To use the Raft implementation in your own project, import the package after running the following command
```bash
go get github.com/amodkala/raft/pkg/raft
```

## Use (Container)
A container image running a demo Raft peer is available via the GitHub Container Registry
```bash
docker pull ghcr.io/amodkala/raft:main
```

## Build (Binary)
With a Flake-enabled Nix installation, simply run
```bash
nix build github:amodkala/raft#bin
```
NOTE: the WAL is written to /var/{id}.wal, so sudo privileges are required
```bash
sudo ./result/bin/cmd
```

## Build (Container Image)
With a Flake-enabled Nix installation:
```bash
nix build github:amodkala/raft#docker
```

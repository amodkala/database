# raft

Implementation of the Raft consensus protocol for distributed log consistency

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

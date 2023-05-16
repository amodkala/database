## note
This project is a work in progress and should be regarded as such

# about
### description
This repository contains the source code for a distributed database which guarantees (eventual) consistency by way of the [Raft consensus algorithm](https://raft.github.io/raft.pdf)

### implementation details
- The database is intended to be run in configurations where cluster size and member network addresses are known (as in Kubernetes StatefulSets or some headless Deployments). This is to make handling cluster membership changes easier.
- Logging and storage are both handled using an LSM-Tree, the code for which is in a [separate repository](https://github.com/amodkala/lsmtree)

Atomicity with Raft
- An entry is applied to the log upon receipt, but only written to storage once
  committed
- We give every transaction a unique ID, attach it to each of its entries
- Central to the atomicity guarantee is the notion that even the successfully
  committed entries of a failed transaction should be scrubbed from the database

Design decisions
- Records stored in LSM Tree need metadata to aid in reconciliation, usually in
  the form of a timestamp. Timestamps should be added when the client first 
  receives the upsert 
- Added timestamp field to Entry type

- After reading about ACID guarantees, I should probably structure writes using
  a transaction type. Each record should have a transactionId for atomicity
  (undo all writes from a partially failed tx)
- Added transaction type for atomicity 

- Now that this module contains the functionality of several connected
  components, I have to think about how to organize code. For example, when the
  client first receives a new record it must be written to the raft leader's log
  and replicated before it can be committed. It is the responsibility of the
  client to report whether a record is successfully committed, it is within the
  purview of the raft consensus module to to do the actual committing, but the
  log is itself used by the LSM tree for fault tolerance.

- Another issue is that the log is a file stored on disk which is frequently
  accessed (reads and writes) by the Raft module, but the lsm tree is also 
  responsible for clearing log segments related to flushed memtables

- Now that I've decided to use protobuf over the wire via gRPC (Raft) and for
  encoding records and writing them to disk (WAL, SSTable), I have to think more
  about whether to define schemas per context, how best to represent
  information, and more.

- How to reconcile Raft's notion of log compaction with the LSM Tree? When
  memtable contents are flushed to disk, standard LSM Tree practice is to delete
  the relevant part of the log. Data compaction occurs within the SSTable
  compaction process, but the WAL is not updated as far as I can tell.

atomicity:
- transactions applied to the raft leader's log (wal Write)
- transactions replicated across cluster (automatically via Raft)
- once replicated, leader commits to storage (how to signal this)
- once leader commits, cluster members commit
    - transaction considered complete at this point

what if CM owns WAL?
- have to expose a public method to database client for replication
- client calls cm.Replicate, passes transaction
    - if leader, cm writes transaction to WAL
    - replication as usual
- LSM Tree doesn't need to read the WAL directly, can just send it the entries
    - this took me like multiple days to figure out -_-

can we solve transaction atomicity with a channel of channels?
- would be a channel of (tx id, commit)

Raft Demo
- handler at /write route expects array of JSON {key, value} objects, passes to
  consensus module
  - if CM isn't leader, handler forwards request to raft leader

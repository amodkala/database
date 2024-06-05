Fundamentals:
- Raft
    - Leader receives transaction from client, appends to log
    - leader log entries are replicated to the cluster, 

Scopes/Responsibilities
- Database client
    - transforms sequence of key/value pairs into a transaction
    - redirects requests if its consensus module is not the leader
    - else issues command to raft consensus module to replicate transaction
    - upon successful 
- Raft
    - leader replicates transaction to cluster
    - leader 
- WAL
    - persist records to disk once transaction is replicated
- Memtable
    - once transaction is committed to log, entries are written to memtable

Assumptions
- Replication can't be modelled as a simple request/response from the client
  perspectivce

Schema per context
- Raft CM processes {term, key, value}
- WAL processes {term, key, value}
- Memtable processes {key, value}
- SSTables process {key, value}
- Bloom filters process {key}

Considerations
- Records stored in LSM Tree need metadata to aid in reconciliation, usually in
  the form of a timestamp. Timestamps should be added when the client first 
  receives the upsert 
- After reading about ACID guarantees, I should probably structure writes using
  a transaction type. Each record should have a transactionId for atomicity
  (undo all writes from a partially failed tx)

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

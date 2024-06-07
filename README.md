# database

This project started off as a toy implementation of the Raft consensus algorithm and promptly turned into an exercise in scope-creep and yak-shaving. I'm now writing a toy *everything*: replication via Raft, write-ahead log for durability, LSM Tree for disk-resident storage, bloom filters to reduce read amplification, and hopefully that's it. 

Designing Data-Intensive Systems by Martin Kleppmann and Database Internals by Alex Petrov have been major influences on my design choices for this project, and I have nothing but positive things to say about either book (except that they've cost me dozens of hours!) I take rough notes on the chapters that I go through, available for your reading pleasure on my website.

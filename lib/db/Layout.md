somaadm DB layout
=================

Bucket jobs
  Bucket active
    Timestamp -> jobid
  Bucket finished
    Timestamp -> jobid
  Bucket data
    jobid -> jobdata

Bucket tokens
  Bucket user
    Expire-Timestamp -> token
  Bucket admin
    Expire-Timestamp -> token

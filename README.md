# Distributed scheduler

Tiny distributed scheduler.
For distributed locks, dlocker package is used.

## Scheme

Execute action during lock period.
For critical data updates (during lock), it is recommended to additionally use data versioning.

For more details see examples.

# GoPass

a rewrite of [https://github.com/JHaals/yopass](https://github.com/JHaals/yopass) in Go.

Requires a local instance of memcached.


Store secrets AES encrypted in memory(memcached) for a fixed period of time.
Secrets can then be shared more securely over channels such as IRC and Email.

* AES-256 encryption
* Secrets can only be viewed once
* No secrets are written to disk
* No accounts and user management required
* Secrets self destruct after X hours

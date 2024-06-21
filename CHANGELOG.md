CHANGELOG
=========

[Unreleased]
-------

- Fixed: distinguish tickets based on `startAt+fingerprint` instead of `fingerprint`.
- Improved: ticket messages
- Improved: matrix client, it doesn't use DB anymore.
- Changed: `!!ticket` command to get a param to specify the ticket format.
- Removed: `!!ticket_yaml` command. use `!!ticket yaml`.

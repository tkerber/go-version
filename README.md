# go-version

go-version creates generates version constants based on RCS data.

The following RCSs generate the given constants:

## Git

* `CommitHashShort` string of the short commit hash.
* `CommitHashLong` string of the full commit hash.
* `CommitTag` string of the last used commit tag (or the empty string).
* `CommitTagIsExact` bool if the tag applies exactly to the latest commit.
* `CommitDate` a `time.Time` struct of the date of the last commit.

## Bazaar

* `RevisionId` string of the full revision id.
* `RevNo` integer of the current revision number.
* `CommitTag` string of the last used commit tag (or the empty string).
* `CommitTagIsExact` bool if the tag applies exactly to the latest commit.
* `CommitDate` a `time.Time` struct of the date of the last commit.

## Mercurial

* `CommitHash` string of the full commit hash.
* `RevNo` integer of the current revision number.
* `CommitTag` string of the last used commit tag.
* `CommitTagIsExact` bool if the tag applies exactly to the latest commit.
* `CommitDate` a `time.Time` struct of the date of the last commit.

`sml` is a command line tool to sync from the Go langauge super repo
at [smallrepo.com][1]. It provides an easy way to manage your Go
language repositories.

[1]: https://smallrepo.com

## Usage

**Step 1**, get and install the client side tool:

```
go get -u shanhu.io/sml
```

**Step 2**, track the repositories that you care, for example:

```
sml track shanhu.io/aries
```

The repository must in the set of repositories that smallrepo tracks.
If you want [smallrepo.com][1] to track more repositories,
[file an issue][2].

[2]: https://github.com/shanhuio/sml/issues/new?title=Track+new+repo

**Step 3**, fetch the latest version of these repos:

```
sml sync
```

It synchronizes to the HEAD of [smallrepo][1], which is guaranteed to
be buildable.

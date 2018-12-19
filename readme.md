`sml` is a command line tool to sync from the Go langauge super repo
at [gopkgs.io][1]. It provides a simple way to manage your Go
language repositories.

[1]: https://gopkgs.io

## How to start?

**Step 1**, get and install the client side tool:

```
go get -u shanhu.io/sml
```

**Step 2**, track the repositories that you care, for example:

```
sml track shanhu.io/aries
```

The repository must in the set of repositories that gopkgs.io tracks.
If you want [gopkgs.io][1] to track more repositories,
[file an issue][2].

[2]: https://github.com/shanhuio/sml/issues/new?title=Track+new+repo

**Step 3**, fetch the latest version of these repos:

```
sml sync
```

It synchronizes to the HEAD of [gopkgs.io][1], which is guaranteed to
be buildable.

## Usage

- `sml track`: prints all tracking repositories.
- `sml track <repo> ...`: add repositories into the tracking set.
- `sml untrack <repo> ...`: remove repositories from the tracking set.
- `sml sync`: merge HEAD from gopkgs.io into the tracking repositories.

## Contact

Email: liulonnie@gmail.com

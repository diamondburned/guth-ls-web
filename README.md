This application requires:

- https://github.com/TeamUlysses/utime
- https://github.com/Guthen/guthlevelsystem

Usage:

0. `go build`
1. Edit `.env` appropriately.
2. Write a systemd service file.
3. Start.

Overriding:

1. In the same directory as the binary, create the following directory tree:

```
- guth-ls-web # binary
- frontend/
	- components/
		- header.html # optional
		- footer.html # optional
```

2. Edit the empty files to add custom branding or anything.
3. Optionally copy `root/index.html` from the source repo to mirror the
   directory layout and edit that as well.

Updating:

1. `git pull` and rebuild.
2. Override the binary.
3. Restart.

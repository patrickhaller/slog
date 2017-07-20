### RATIONALE 

Dave Cheney opined that there were only two kinds of logs needed:
those used by production admin staff to keep the app running, and 
development debug logs.
[^1](https://dave.cheney.net/2015/11/05/lets-talk-about-logging)

That seems correct; it just needs audit / accounting logs to cover 
the vast majority of use cases. 

Second, production logs should be keyed, where the key 1) indicates the
locality of the problem, and 2) is relatively persistent across releases.
By default, slog will generate a key from the current source filename 
and function name for production logs. 

This is for everyone who has googled log output, only to find that it was specific to 
only a few releases, and all the results are just the code itself.

### CONFIG

The config has 4 options: 
```
	File = the filename for production or development logs 
	Debug = true for dev logs, false for prod
	Prefix = keyword prefix for the locality hash 
	AuditFile = filename for audit logs -- optional
```


### EXAMPLES

A barebones init:

```golang
	slog.Init(slog.Config{
		File:      "/dev/stderr",
		Debug:     false,
		Prefix:    "TEST",
	})
	slog.D("open file `%s'", cfg.File)
	//...
	slog.P("failed accessing `%s': %v", cfg.File, err)
```

Which would yield the following because Config.Debug was false:

```
2017/07/19 21:07:44 WARN TEST-B0EAE04C failed accessing `/tmp/a': File not found
```

While a maximally verbose init:

```golang
	slog.Init(slog.Config{
		File:      "/dev/stderr",
		Debug:     true,
		AuditFile: "/var/log/TEST.log",
		Prefix:    "TEST",
	})
	slog.D("open file `%s'", cfg.File)
	//...
	slog.P("failed accessing `%s': %v", cfg.File, err)
	//...
	slog.A("%s `%s' %s", r.Method, r.URL, r.RemoteAddr)
```

```
2017/07/19 21:33:14.069128 go-log.go:53: open file `/tmp/a'
2017/07/19 21:33:14.069179 go-log.go:73: WARN TEST-B0EAE04C failed accessing `/tmp/a': File not found
```



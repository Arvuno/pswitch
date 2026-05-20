# Logging

`pswitch` uses global helper functions:

- `logx.Infof(...)`
- `logx.Warnf(...)`
- `logx.Debugf(...)`

Color handling:

- auto on when the terminal supports color
- auto off when it does not
- manual override with `--log-color=true|false`

Examples:

```bash
./bin/pswitch serve --log-color=true
./bin/pswitch serve --log-color=false
```

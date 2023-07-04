# nglog - Formats (*php-fpm* + *nginx*) logs

**nglog**  is a cli program, which splits and formats typical PHP error messages from a nginx error log file.


### Without:
```shell
tail /var/log/nginx/example.com_error.log

2023/07/03 03:03:37 [error] 396782#396782: *300397 FastCGI sent in stderr: "PHP message: PHP Warning:  Undefined array key "key0" in /var/www/example.com/index.php on line 12PHP message: PHP Warning:  Undefined array key "key1" in /var/www/example.com/index.php on line 13" while reading upstream, client: 127.0.0.1, server: example.com, request: "GET / HTTP/2.0", upstream: "fastcgi://unix:/run/php/php8.0-fpm.sock:", host: "example.com"
```

### With nglog:
```shell
tail /var/log/nginx/example.com_error.log | nglog

2023/07/03 03:03:37 [error] 396782#396782: *300397 FastCGI sent in stderr: "PHP Warning:  Undefined array key "key0" in /var/www/example.com/index.php on line 12" while reading upstream, client: 127.0.0.1, server: example.com, request: "GET / HTTP/2.0", upstream: "fastcgi://unix:/run/php/php8.0-fpm.sock:", host: "example.com"
2023/07/03 03:03:37 [error] 396782#396782: *300397 FastCGI sent in stderr: "PHP Warning:  Undefined array key "key1" in /var/www/example.com/index.php on line 13" while reading upstream, client: 127.0.0.1, server: example.com, request: "GET / HTTP/2.0", upstream: "fastcgi://unix:/run/php/php8.0-fpm.sock:", host: "example.com"
```

```shell
tail /var/log/nginx/example.com_error.log | nglog -t "%ts% - %php%"

2023/07/03 03:03:37 - PHP Warning:  Undefined array key "key0" in /var/www/example.com/index.php on line 12
2023/07/03 03:03:37 - PHP Warning:  Undefined array key "key1" in /var/www/example.com/index.php on line 13
```


# Install

### x64 (linux-amd64)
```shell
sudo wget -q -O /bin/nglog https://github.com/marcelhencke/nglog/releases/latest/download/nglog-amd64; sudo chmod +x /bin/nglog
```

### arm64 (linux-arm64)
```shell
sudo wget -q -O /bin/nglog https://github.com/marcelhencke/nglog/releases/latest/download/nglog-arm64; sudo chmod +x /bin/nglog
```



# Usage
```shell
# log file as argument:
nglog /var/log/nginx/example.com_error.log

# tail reads log file and transfer data via a pipe to nglog:
tail -f /var/log/nginx/example.com_error.log | nglog

# custom template for log file:
nglog -t "%ts% - %php% - %ng_upstream%" /var/log/nginx/example.com_error.log
```
### Template Keys

| Key       | Description                                | Example                                                                                                                                                                                                                                                                                                                                 |
|-----------|--------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| raw       | no further line manipulation               | 2023/07/03 03:03:37 [error] 396782#396782: *300397 FastCGI sent in stderr: "PHP Warning:  Undefined array key "key0" in /var/www/example.com/index.php on line 12" while reading upstream, client: 127.0.0.1, server: example.com, request: "GET / HTTP/2.0", upstream: "fastcgi://unix:/run/php/php8.0-fpm.sock:", host: "example.com" |
| prefix    | first match group from LogLineFastCGIRegex | 2023/07/03 03:03:37 [error] 396782#396782: *300397 FastCGI sent in stderr: "                                                                                                                                                                                                                                                            |
| suffix    | third match group from LogLineFastCGIRegex | " while reading upstream, client: 127.0.0.1, server: example.com, request: "GET / HTTP/2.0", upstream: "fastcgi://unix:/run/php/php8.0-fpm.sock:", host: "example.com"                                                                                                                                                                  |
| ts        | timestamp                                  | 2023/07/03 03:03:37                                                                                                                                                                                                                                                                                                                     |
| php       | PHP message                                | PHP Warning:  Undefined array key "key0" in /var/www/example.com/index.php on line 12                                                                                                                                                                                                                                                   |
| ng_*xxxx* | nginx var, e.g. ng_server, ng_upstream ... | example.com, "fastcgi://unix:/run/php/php8.0-fpm.sock:"                                                                                                                                                                                                                                                                                 |


### Debug Flags
| Flag            | Default | Description                                                                          |
|-----------------|---------|--------------------------------------------------------------------------------------|
| -d --debugMode  | false   | Print out debug logs. Helpful when defining custom regex.                            |
| --debugMaxLines | 0       | Only read x lines of input data. Only effective when debugMode is enabled and x > 0. |


### Overwrite Flags

This flag overwrite core regex's and strings, which nglog uses to split and format incoming log lines. Customize these in case your log lines differ from mine. 

| Flag                            | Default                                                                                                                                            | Description                                                                                                                                                                                                                                                                                                               |
|---------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| --overwriteLogLineCompleteRegex | `" while reading[A-z ]* upstream`                                                                                                                  | This RegEx tests whether the lines read are a complete log line.                                                                                                                                                                                                                                                          |
| --overwriteLogLineFastCGIRegex  | `^(\d{4}[\/-]\d{2}[\/-]\d{2} \d{2}:\d{2}:\d{2} \[\w+\] .+ FastCGI sent in stderr: ")([\s\S]*)(" while reading[A-z ]* upstream(?:, \w+: "?.+"?)*)$` | This RegEx tests whether the log line is a FastCGI log line.<br/>This expression are divided in 3 matching groups: prefix, body, suffix.<br/><br/>**Prefix:** Timestamp and FastCGI identifier.<br/>**Body:** Contains the PHP messages.<br/>**Suffix:** Some nginx properties, like client, server, upstream, host, etc. |
| --overwriteLogLinePhpMsgSplit   | `PHP message: `                                                                                                                                    | This string is the split value for PHP messages.                                                                                                                                                                                                                                                                          |
| --overwriteNginxVarRegex        | `, (\w+): ("?[^,]+"?)`                                                                                                                             | This RegEx finds the nginx vars in the line suffix.                                                                                                                                                                                                                                                                       |



# Credits
https://github.com/napicella/go-linux-pipes

https://github.com/spf13/cobra

# License
nglog is released under the Apache 2.0 license. See [LICENSE](LICENSE)

# nglog - Formats (*php-fpm* + *nginx*) logs

**nglog**  is a cli program, which splits and formats typical PHP error messages from a nginx error log file.


### Without:
```
user@server:~# tail /var/log/nginx/example.com_error.log

2023/07/03 03:03:37 [error] 396782#396782: *300397 FastCGI sent in stderr: "PHP message: PHP Warning:  Undefined array key "key0" in /var/www/example.com/index.php on line 12PHP message: PHP Warning:  Undefined array key "key1" in /var/www/example.com/index.php on line 13" while reading upstream, client: 127.0.0.1, server: example.com, request: "GET / HTTP/2.0", upstream: "fastcgi://unix:/run/php/php8.0-fpm.sock:", host: "example.com"
```

### With nglog:
```
user@server:~# tail /var/log/nginx/example.com_error.log | nglog

2023/07/03 03:03:37 [error] 396782#396782: *300397 FastCGI sent in stderr: "PHP Warning:  Undefined array key "key0" in /var/www/example.com/index.php on line 12" while reading upstream, client: 127.0.0.1, server: example.com, request: "GET / HTTP/2.0", upstream: "fastcgi://unix:/run/php/php8.0-fpm.sock:", host: "example.com"
2023/07/03 03:03:37 [error] 396782#396782: *300397 FastCGI sent in stderr: "PHP Warning:  Undefined array key "key1" in /var/www/example.com/index.php on line 13" while reading upstream, client: 127.0.0.1, server: example.com, request: "GET / HTTP/2.0", upstream: "fastcgi://unix:/run/php/php8.0-fpm.sock:", host: "example.com"
```

```
user@server:~# tail /var/log/nginx/example.com_error.log | nglog -t "%ts% - %php%"

2023/07/03 03:03:37 - PHP Warning:  Undefined array key "key0" in /var/www/example.com/index.php on line 12
2023/07/03 03:03:37 - PHP Warning:  Undefined array key "key1" in /var/www/example.com/index.php on line 13
```



# Usage
```
log file as argument:
user@server:~# nglog /var/log/nginx/example.com_error.log

tail reads log file and transfer data via a pipe to nglog:
user@server:~# tail -f /var/log/nginx/example.com_error.log | nglog

custom template for log file:
user@server:~# nglog -t "%ts% - %php% - %ng_upstream%" /var/log/nginx/example.com_error.log
```


[//]: # (TODO: all params)


# Credits
https://github.com/napicella/go-linux-pipes
https://github.com/spf13/cobra

# License
nglog is released under the Apache 2.0 license. See [LICENSE](LICENSE)

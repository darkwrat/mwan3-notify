# mwan3-notify

Get yourself a notification straight to laptop when some mwan3-managed network is down.

/etc/mwan3.user:
```
#!/bin/sh

SECRET="xxx"

/usr/bin/curl --insecure \
	--data-urlencode "hostname=${HOSTNAME}" \
	--data-urlencode "action=${ACTION}" \
	--data-urlencode "interface=${INTERFACE}" \
	--data-urlencode "device=${DEVICE}" \
	--data-urlencode "secret=${SECRET}" \
	"https://mwan3-notify-addr.a/mwan3-notify" \
	-o /dev/null >/dev/null 2>&1
```
nginx.conf:
```
        location /mwan3-notify {
            fastcgi_pass unix:/var/run/mwan3-notify-fcgi/fcgi.sock;
            include fastcgi_params;
        }
```
command line:
```
bin/mwan3-notify-fcgi -s xxx -i /usr/share/icons/gnome/32x32/emblems/emblem-new.png
```

# mwan3-notify

Get yourself a notification straight to the laptop when some mwan3-managed network is down.

![single](https://github.com/darkwrat/mwan3-notify/raw/master/doc/1_single.png)
![multiple](https://github.com/darkwrat/mwan3-notify/raw/master/doc/2_multiple.png)

/etc/mwan3.user (cgi-bin/luci/admin/network/mwan/notify):
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
nginx.conf for 443 (and reload):
```
        location /mwan3-notify {
            fastcgi_pass unix:/run/mwan3-notify-fcgi/fcgi.sock;
            include fastcgi_params;
        }
```
/etc/tmpfiles.d/mwan3-notify.conf:
```
d /run/mwan3-notify-fcgi 0755 <your-user> <your-group> -
```
create the sock dir:
```
sudo systemd-tmpfiles --create
```
make and install the binary:
```
make
sudo cp bin/mwan3-notify-fcgi /usr/local/bin
```
test by hand:
```
/usr/local/bin/mwan3-notify-fcgi -s xxx -i /usr/share/icons/gnome/32x32/emblems/emblem-new.png &
curl --insecure --data-urlencode "hostname=a" --data-urlencode "action=b" --data-urlencode "interface=c" --data-urlencode "device=d" --data-urlencode "secret=xxx" "https://127.0.0.1/mwan3-notify"
fg
```
add to autostart ~/.config/autostart/mwan3-notify.desktop:
```
[Desktop Entry]
Type=Application
Hidden=false
X-GNOME-Autostart-enabled=true
Exec=/usr/local/bin/mwan3-notify-fcgi -s xxx -i /usr/share/icons/gnome/32x32/emblems/emblem-new.png
Comment=mwan3-notify
```
relogin and forget about it.

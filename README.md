HubbIT-Sniffer
==================

This is a simple python script that uses [tshark](https://www.wireshark.org/docs/man-pages/tshark.html) to capture mac addresses and send them to [hubbIT](https://github.com/cthit/hubbIT). 


***Note** this script requires ```python 3.x```, ```requests``` [1] and a wireless interface that supports *monitor* and *promiscuous* mode. 

[1]: http://www.python-requests.org/en/latest/        “Python Requests”
## How to use:
``` $ cd /path/to/hubbIT-sniffer ```

First you need to setup the required interface. Fortunately there is script provided that will do just that:
```
# scripts/setup-hubbit-iface.sh
```

After that you only need to start the *hubbit-sniffer*:
```
$ screen -S hubbit-sniffer -U
$ python3 sniffer/sniffer.py -a api_key_goes_here -b sniffer/blacklist.txt
```

If you need more customisability, have a look at ```python3 sniffer/sniffer.py --help```

### Blacklist
The blacklist is a ```\n``` separated text file of mac addresses to ignore. This is best used together with the script ```nearby-aps.sh``` in order to ignore any traffic coming from nearby access points.

```
# scripts/nearby-aps.sh > sniffer/blacklist.txt
```

### Timeout
Using the ```-t``` flag you can adjust how often the sniffer sends batch updates to the server. Default is 5 seconds.

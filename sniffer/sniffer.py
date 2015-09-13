# This is the application that sends MAC-addresses to the hubbIT server
#!/usr/bin/python3

import reqs
import blacklist
import subprocess
import threading
import time
import signal
import requests
import argparse

class MacStorage():
    def __init__(self):
        self._lock = threading.Lock()
        self._macs = dict()

    def seen(self, mac):
        ''' Add mac to dict, where time is the time the mac was last seen '''
        self._lock.acquire()
        self._macs[mac] = self._macs.get(mac, 0) + 1
        self._lock.release()

    def clear(self):
        self._lock.acquire()
        self._macs.clear()
        self._lock.release()

    def list_and_clear(self):
        self._lock.acquire()
        maclist = [(k,v) for k,v in self._macs.items()]
        if not maclist:
            maclist.append("")
        self._macs.clear()
        self._lock.release()

        return maclist

class Capture(threading.Thread):
    def __init__(self, storage, iface, blacklist_path):
        self._storage = storage
        self._iface = iface
        self._blacklist_path = blacklist_path
        threading.Thread.__init__(self)

    def _build_command(self):
        filter = blacklist.Blacklist(self._blacklist_path).create_filter()
        return ['tshark', '-i', self._iface , '-p', '-l', '-n', '-T', 'fields', '-e', 'wlan.sa', filter]

    def run(self):
        self._tshark_proc = subprocess.Popen(self._build_command(), stdout=subprocess.PIPE, stderr=subprocess.DEVNULL)
        lines_it = iter(self._tshark_proc.stdout.readline, b'')
        for sa in lines_it:
            if len(sa) > 1:
                self._storage.seen(sa.decode('utf-8').strip())

    def stop(self):
        self._tshark_proc.terminate() # Send sigterm to subprocess


class Main:
    def __init__(self, api=None, url=None, interface=None, blacklist_path=None, timeout=5):
        self._storage = MacStorage()
        self._keep_capturing = True
        self._sigint = False
        self._url = url
        self._api = api
        self._iface = interface
        self._blacklist_path = blacklist_path
        self._timeout = timeout

    def handle_sigusr1(self, signal, frame):
        print("Caught SIGUSR1, reloading blacklist")

        if self._cap is not None:
            self._cap.stop()
            self._keep_capturing = False

    def handle_sigint(self, signal, frame):
        print("Caught SIGINT")
        self._cap.stop()
        self._keep_capturing = False
        self._sigint = True

    def PUT_to_server(self, payload):
        r = requests.put(self._url,
                         headers={"Authorization":"Token token=" + self._api},
                         json=payload)
        return r.status_code, r.reason


    def run(self):
        signal.signal(signal.SIGUSR1, self.handle_sigusr1)
        signal.signal(signal.SIGINT, self.handle_sigint)

        while not self._sigint:
            self._keep_capturing = True
            self._cap = Capture(self._storage, self._iface, self._blacklist_path)
            self._cap.start()
            while self._keep_capturing:
                time.sleep(self._timeout)
                macs = self._storage.list_and_clear()
                status_code, reason = self.PUT_to_server({"macs":macs})
                print(time.strftime("%F %T") + " - " + str(len(macs)) + " -> " + self._url  + " <-- " + str(status_code) + " " + reason)
            self._storage.clear()
            self._cap.join()

def main():
    parser = argparse.ArgumentParser(description="Mac sniffer")
    parser.add_argument('-a', '--api', default="PLZ_LET_ME_IN", help="The API-key to send to the server")
    parser.add_argument('-u', '--url', default="https://hubbit.chalmers.it/sessions.json", help="Full PUT address to the server")
    parser.add_argument('-i', '--interface', metavar="IFACE", default="hubbit", help="The monitor interface. Note: must be in promiscious mode")
    parser.add_argument('-b', '--blacklist', metavar="PATH", default="blacklist.txt", help="The path to the blacklist file.")
    parser.add_argument('-t', '--timeout', metavar="SECONDS", type=int, default=5, help="Time in seconds between each batch upload of macs to the server")

    args = parser.parse_args()

    req = reqs.Requirements(iface=args.interface)
    if not req.check():
        print('Fix the above before continuing')
        exit(1)

    m = Main(api=args.api, url=args.url, interface=args.interface,
             blacklist_path=args.blacklist, timeout=args.timeout)
    m.run()

if __name__ == '__main__':
    main()

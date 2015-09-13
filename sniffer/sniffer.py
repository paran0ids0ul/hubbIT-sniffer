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
        self._macs.clear()
        self._lock.release()

        return maclist

class Capture(threading.Thread):
    def __init__(self, storage, iface):
        self._storage = storage
        self._iface = iface
        threading.Thread.__init__(self)

    def _build_command(self):
        filter = blacklist.Blacklist().create_filter()
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
    def __init__(self, api=None, url=None, interface=None):
        self._storage = MacStorage()
        self._keep_capturing = True
        self._sigint = False
        self._url = url
        self._api = api
        self._iface = interface

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


    def run(self):
        signal.signal(signal.SIGUSR1, self.handle_sigusr1)
        signal.signal(signal.SIGINT, self.handle_sigint)

        while not self._sigint:
            self._keep_capturing = True
            self._cap = Capture(self._storage, self._iface)
            self._cap.start()
            while self._keep_capturing:
                time.sleep(5)
                macs = self._storage.list_and_clear()
                pload = {"macs":macs}
                r = requests.put(self._url,
                                 headers={"Authorization":"Token token=" + self._api},
                                 json=pload)
                print(time.strftime("%F %T") + " - " + str(len(macs)) + " -> " + self._url  + " <-- " + str(r.status_code))
            self._storage.clear()
            self._cap.join()

def main():
    parser = argparse.ArgumentParser(description="Mac sniffer")
    parser.add_argument('-a', '--api', default="PLZ_LET_ME_IN", help="The API-key to send to the server")
    parser.add_argument('-u', '--url', default="https://hubbit.chalmers.it/sessions.json", help="Full PUT address to the server")
    parser.add_argument('-i', '--interface', metavar="IFACE", default="hubbit", help="The monitor interface. Note: must be in promiscious mode")

    args = parser.parse_args()

    req = reqs.Requirements(iface=args.interface)
    if not req.check():
        print('Fix the above before continuing')
        exit(1)

    m = Main(api=args.api, url=args.url, interface=args.interface)
    m.run()

if __name__ == '__main__':
    main()

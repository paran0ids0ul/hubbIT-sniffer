# This is the application that sends MAC-addresses to the hubbIT server
#!/usr/bin/python3

import reqs
import blacklist
import subprocess
import threading
import time
import signal

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
    def __init__(self, storage):
        self._storage = storage
        threading.Thread.__init__(self)

    def _build_command(self):
        filter = blacklist.Blacklist().create_filter()
        return ['tshark', '-i' 'hubbit', '-p', '-l', '-n', '-T', 'fields', '-e', 'wlan.sa', filter]

    def run(self):
        self._tshark_proc = subprocess.Popen(self._build_command(), stdout=subprocess.PIPE, stderr=subprocess.DEVNULL)
        lines_it = iter(self._tshark_proc.stdout.readline, b'')
        for sa in lines_it:
            if len(sa) > 1:
                self._storage.seen(sa.decode('utf-8').strip())

    def stop(self):
        self._tshark_proc.terminate() # Send sigterm to subprocess


class Main:
    def __init__(self):
        self._storage = MacStorage()
        self._keep_capturing = True
        self._sigint = False

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
            self._cap = Capture(self._storage)
            self._cap.start()
            while self._keep_capturing:
                time.sleep(5)
                # Instead of print, PUT that shiet to the server
                print(self._storage.list_and_clear())
            self._storage.clear()
            self._cap.join()

def main():
    req = reqs.Requirements()
    if not req.check():
        print('Fix the above before continuing')
        exit(1)
    m = Main()
    m.run()

if __name__ == '__main__':
    main()

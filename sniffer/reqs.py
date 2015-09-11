from shutil import which
from subprocess import Popen

class Requirements():
    """This class manages and checks requirements before running the sniffer"""

    def __init__(self, iface='mon0'):
        self.iface = iface

    def check(self):
        have_reqs = True
        if not which('ip'):
            print('Unable to find ip program. Unable to change to promiscious mode')
            have_reqs = False
        if not which('iw'):
            print('Unable to find iw. Needed for checking wireless interface')
            have_reqs = False

        if not self._mon_mode():
            print("Monitor mode not enabled")
            have_reqs = False
        return have_reqs

    def _mon_mode(self):
        print("Hello")
        with Popen(["iw", "dev", self.iface, "info"], shell=True) as iw:
            print(iw.stdout.read())
        return False

    def _promisc(self):
        return False




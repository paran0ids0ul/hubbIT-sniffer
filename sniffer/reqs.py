class Requirements():
    """This class manages and checks requirements before running the sniffer"""

    def __init__(self, iface='hubbit'):
        self.iface = iface

    def check(self):
        have_reqs = True
        if not self._iface_exists():
            print("Network interface \'" + self.iface +
                  "\' doesn't exist. Run scripts/setup-mon-iface.sh as root")
            have_reqs = False
        return have_reqs

    def _iface_exists(self):
        content = None
        with open("/proc/net/dev") as f:
            content = f.read()

        return self.iface in content

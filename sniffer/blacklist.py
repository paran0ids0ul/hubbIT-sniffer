from enum import Enum

class Filter(Enum):
    BPF = 1
    Capture = 2
    Display = 3

class Blacklist:
    def __init__(self):
        self.blacklist_file = "blacklist.txt"
        self.default_blacklist = ["00:1f:9d:b6:e0:00", "00:00:00:00:00:00", "ff:ff:ff:ff:ff:ff", "2c:54:2d:3a:5f:60"]

    def _read_blacklist(self):
        content = None
        try:
            with open(self.blacklist_file) as f:
                content = f.readlines()
        except IOError:
            return []
        return [x.strip().lower() for x in content]


    def create_filter(self, type=Filter.Capture):
        macs = self._read_blacklist()
        macs.extend(self.default_blacklist)
        filter = "wlan.sa !="
        if type == Filter.Capture or type == Filter.BPF:
            filter = "not wlan src"

        if len(macs) > 0:
            return filter + " " + (" and " + filter + " ").join(macs)
        return ""



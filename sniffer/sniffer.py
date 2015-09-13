# This is the application that sends MAC-addresses to the hubbIT server
#!/usr/bin/python3

import reqs
import blacklist
import pyshark

def main():
    blist = blacklist.Blacklist()
    #req = reqs.Requirements()
    #if not req.check():
    #    print('Fix the above before continuing')
    #    exit(1)
    cap = pyshark.LiveCapture(interface="hubbit", only_summaries=False, bpf_filter=blist.create_filter(blacklist.Filter.Capture))
    for frame in cap.sniff_continuously():
        wlan = frame['wlan']
        if hasattr(wlan, "sa"):
            print(wlan.sa)

if __name__ == '__main__':
    main()

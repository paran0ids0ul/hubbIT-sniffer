# This is the application that sends MAC-addresses to the hubbIT server
#!/usr/bin/python3

import reqs

def main():
    req = reqs.Requirements()
    if not req.check():
        print('Fix the above before continuing')
        exit(1)


if __name__ == '__main__':
    main()

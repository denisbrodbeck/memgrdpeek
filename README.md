# memgrdpeek

## Intention

memgrdpeek tries to evaluate the effective protection of the golang library [memguard](https://github.com/awnumar/memguard).

The go app uses a `memguard.LockedBuffer` to protect its data. A [python script](memory_reader.py) tries to extract protected memory areas from the running go app. This script reads all used memory areas from `/proc/app-pid/maps` and tries to find 32 byte long areas within `/proc/app-pid/mem`.

## How to use

Tested on **Linux** only! Might work on *OS X*, but not on *Windows*.

Get this project:

```bash
go get github.com/denisbrodbeck/memgrdpeek
cd "$GOPATH/src/github.com/denisbrodbeck/memgrdpeek"
```

* [build using `make`]
* run `memgrdpeek`
* open another terminal
* run `sudo memory_reader.py`

## What it does

Essentially we try to create a data structures with a memory protected area within it:

```golang
type message struct {
  ID        []byte                 // 32 byte
  Recipient []byte                 // 32 byte
  Sender    []byte                 // 32 byte
  Message   *memguard.LockedBuffer // Gets 32 byte input
  Meta      []byte                 // 32 byte
}
```

Create a new message:

```golang
msg := newMessage(
  []byte("(((id    | 1234567890123456  )))"),
  []byte("(((recip | bobby@spacer.com  )))"),
  []byte("(((sender| eva@secretmoon.m  )))"),
  []byte("(((msg   | VERY SECRET VERY  )))"), // this should be protected
  []byte("(((meta  | 2017-10-21 10:26  )))"),
  )
```

After reading the process's memory, one expects to discover, that all buffers are visible except the `Message` buffer.

## Output

Starting memguarded app:

```bash
$ bin/memgrdpeek

Running under PID: 25921
Press enter after running 'memory_reader.py' to exit...

```

And run `memory_reader.py` with the provided `pid`:

```bash
$ sudo ./memory_reader.py --pid 25921

INFO:root:Trying pid 25921
INFO:root:Scanning pid 25921 range 7ffec4585000-7ffec45a6000 [stack];
INFO:root:Found 0 candidate pointers
INFO:root:Scanning pid 25921 range 7f038c647000-7f038c6e7000 ;
INFO:root:Found 0 candidate pointers
INFO:root:Scanning pid 25921 range c41fff8000-c420100000 ;
INFO:root:Found 7 candidate pointers
842350543008:    QT_QPA_PLATFORMTHEME=appmenu-qt5LC_TIME=de_DE.UTF-8SHELL=/usr/bin/zshLC_NUMERIC=de_DE.UTF-8LC_PAPER=de_DE.UTF-8
842350543296:    XAUTHORITY=/home/dev/.XauthorityXDG_SESSION_DESKTOP=ubuntuGDMSESSION=ubuntuUPSTART_EVENTS=xsession startedLC_MEASUREMENT=de_DE.UTF-8
842350544192:    (((id    | 1234567890123456  )))(((recip | bobby@spacer.com  )))(((sender| eva@secretmoon.m  )))(((meta  | 2017-10-21 10:26  )))
842350544224:    (((recip | bobby@spacer.com  )))(((sender| eva@secretmoon.m  )))(((meta  | 2017-10-21 10:26  )))
842350544256:    (((sender| eva@secretmoon.m  )))(((meta  | 2017-10-21 10:26  )))
842350544288:    (((meta  | 2017-10-21 10:26  )))
842350730936:    ` �� ��E �`E �@E ���KPCO��K� ���KpCO` ���J�_d� ��K
INFO:root:Scanning pid 25921 range c000000000-c000001000 ;
INFO:root:Found 0 candidate pointers
INFO:root:Scanning pid 25921 range 56a000-58c000 ;
INFO:root:Found 0 candidate pointers
```

Studding the output one can see, that the message fields are clearly visible from memory. Ouput only skips the memory are of the `Message` field.

> `(((id    | 1234567890123456  )))(((recip | bobby@spacer.com  )))(((sender| eva@secretmoon.m  )))(((meta  | 2017-10-21 10:26  )))`

## Conclusion

Using the described approach I wasn't able to extract the content of the protected buffer.

An experienced low-level/assembly hacker would probably be able to extract more meaningful content. PR's are welcome :)

## Credits

The python extraction code wouldn't be possible without prior work of [Vayu](github.com/Vayu/vault_recover). Vayu's project extracts the unencrypted private key from hashicorps [vault](https://github.com/hashicorp/vault).

## License

The MIT License (MIT) — [Denis Brodbeck](https://github.com/denisbrodbeck). Please have a look at the [LICENSE](LICENSE.md) for more details.

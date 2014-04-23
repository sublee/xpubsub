zmqproxy
========

A ØMQ proxy implementation in Go.

    $ go get github.com/sublee/zmqproxy
    $ go install github.com/sublee/zmqproxy
    $ zmqproxy --help
    Usage of zmqproxy:
        Runs ZeroMQ proxy.
    Options:
      -Q       --queue          choose 'queue' device
      -F       --forwarder      choose 'forwarder' device
      -S       --streamer       choose 'streamer' device
      -f 5561  --frontend=5561  listening port the frontend socket binds to.
      -b 5562  --backend=5562   listening port the backend socket binds to.
               --no-traffic     disable traffic reporting.
               --help           show usage message
    $ zmqproxy -F
    2014/04/24 08:40:47 ZeroMQ 'forwarder' device chosen
    2014/04/24 08:40:47 Traffic reporting enabled
    2014/04/24 08:40:47 Proxying between 5561[XSUB] and 5562[XPUB]...

# plan

    d8/domain       Domain name and registrar parsing
    d8/wire         DNS message parsing
    d8/client       Async UDP4 DNS client
    d8/shell        Interactive console and query tree handling
    d8/cmds         Common query logic
    d8/bin/d8       Program that can fire single queries, back queries, 
                    or listen to TCP/HTTP input
    d8/bin/d8cesr   Crawler that works in UCSD crawler infrastructure

# todo

- message printing
- shell
- cmds

# console design

    $ q liulonnie.net @74.220.195.131
    #1090 auth
    ques liulonnie.net
    answ {
        liulonnie.net a 66.147.240.181 4h
    }
    auth {
        ...
    }
    addi {
        ...
    }
    (in 20ms)

    $ recur liulonnie.net
    // some comment
    q liulonnie.net a @192.228.79.201 {
        #1090 auth
        ques liulonnie.net
        answ {
            liulonnie.net a 66.147.240.181 4h
        }
        auth {
            ...
        }
        addi {
            ...
        }
        (in 20ms)
    }
    // some more comment

    $ ip liulonnie.net
    $ about liulonnie.net


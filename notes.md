# plan

    d8/domain       Domain name and registrar parsing
    d8/wire         DNS message parsing
    d8/client       Async UDP4 DNS client
    d8/qtree        Query tree
    d8/qlogic       Common query logic
    d8/bin/d8       Program that can fire single queries, back queries, 
                    or listen to TCP/HTTP input
    d8/bin/d8cesr   Crawler that works in UCSD crawler infrastructure

# todo

message printing

    $ dig liulonnie.net @74.220.195.131
    q liulonnie.net a @192.228.79.201 
    a +20ms {
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
    }

    $ recur liulonnie.net
    recur liulonnie.net a {
        q liulonnie.net a @192.228.79.201 
        a +20ms {
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
        }
    }


    % recur liulonnie.net ns
    % ips liulonnie.net

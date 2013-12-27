# plan

    d8/domain       Domain name and registrar parsing
    d8/packet       DNS packet parsing
    d8/client       Async UDP4 DNS client
    d8/term         Interactive and recursive terminal for a task
    d8/tasks        Common query logics
    d8/bin/d8       Program that can fire single queries, back queries, 
                    or listen to TCP/HTTP input
    d8/bin/d8cesr   Crawler that works in UCSD crawler infrastructure

# todo

- better handling on message unpack error (especially when id is ...)
- mark offline zone servers (save last unreachable time)
- save in sqlite3 and other output devices

# notes

- sqlite3 is good for secondary dumping but not good for results
- what you are looking for is a zip file
    
    input
    log/
        000_www.google.com
        001_xxx
        002_xxx
        003_xyz
    out/
        000_www.google.com
        001_xxx
        002_yyy
        003_xyz
    err


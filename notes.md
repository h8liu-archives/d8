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

- handling message unpack error properly (especially when id is ...)
- mark offline zone servers (mark last unreachable time)
- concurrent crawling
- save in sqlite3

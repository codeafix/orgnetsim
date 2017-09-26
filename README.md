# orgnetsim
A simulator for Organisational Networks

Send msg to random friend
- blocked friend already has msg
   - check own msgs
      - no msg: send msg to next friend
      - msg: send confirm
         - connected
         - reject all others
- not-blocked msg gets through
   - check own msgs
      - no msg: ?
      - msg: is confirm?
         - connected
         - reject all others
      - msg: is req
         - add to req list
      - msg: is reject
         - look in req list
            - items: send confirm
            - empty: send msg to next friend

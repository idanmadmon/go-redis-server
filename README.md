# go-redis-server
- Follow the drawio design
- change parse to know the type (like integer), also to something more readable like struct request that has command as type (since it has to be)
- export parser and endcoder to package RESP
- know types
- stop from outside and listen to ctrl+c from outside
- create Server struct that has Start() Stop()
- do workers and send cfg to needed worker (commands)
- think if using cobra is necessary
- change from logrus to zap
- read mati's working with errors

- If there is time add atomic set
- Add hash keys and change the the values to pair for faster search
- Convert to working with actor model so you can add another parser from config (parser_worker = 2)

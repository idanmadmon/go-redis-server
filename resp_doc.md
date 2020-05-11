## RESP Commands ##

#### ping: ####
- *1\r\n
- $4\r\nping\r\n

#### set: ####
- *3\r\n
- $3\r\nset\r\n
- $3\r\nkey\r\n
- $5\r\nvalue\r\n

#### setnx: ####
- *3\r\n
- $5\r\nsetnx\r\n
- $5\r\nkeynx\r\n
- $7\r\nvaluenx\r\n

 ####get: ####
- *\r\n
- $3\r\nget\r\n
- $3\r\nkey\r\n


## RESP Replies: ##

#### ping: ####
- +PONG\r\n

#### set: ####
- +OK\r\n

#### setnx: ####
- :0\r\n

#### get: ####
- $5\r\nvalue\r\n

#### get (not found): ####
- $-1\r\n
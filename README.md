# Micro consumer

This project is to demonstrate a simple solution to consume micro blogs and persist to some backend.
You can create and include as many external sources through the producers, currently twitter, and change/add/extend backend storage through consumers. The current sql consumer can scale to more than one instance (highly depended on the backend solution) while on the producer side only one instance per source. All produced data will end in Redis before consumed.

The current model is schema-less, meaning any external json data should more or less go through.

To set up and test locally please install golang and docker runtime. 
Under the folder external you will see a simple setup for redis and mariadb. 
When both containers are running you need to edit the .env.local and change the corresponding values, e.g. twitter bearer token TWITTER_BEARER_TOKEN.  

TODO:
- add tests
- complete docker support including code base
- separate producer/consumer in main file so you can run them as separate instances in separate containers

Example output from program

    bash-5.1$ go run main.go 
    2022/10/16 16:05:02 Starting producers and consumers
    2022/10/16 16:05:12 in consumer cgriwwietx redis length:0
    2022/10/16 16:05:13 Producer ID iqeaoeywrw => NewestID:1581644123396018177 OldestID:1581556016776237056 NextToken:b26v89c19zqg8o3fpzejv6t8flbam47pf4qnyvu2w096l ResultCount:10
    2022/10/16 16:05:13 Producer ID iqeaoeywrw pushed twitter id:1581644123396018177
    2022/10/16 16:05:13 Producer ID iqeaoeywrw pushed twitter id:1581630095462400000
    ...
    2022/10/16 16:05:17 Producer ID iqeaoeywrw pushed twitter id:1579199420012580864
    2022/10/16 16:05:17 Producer ID iqeaoeywrw => NewestID:1579198649024000000 OldestID:1579198600562647040 NextToken: ResultCount:2
    2022/10/16 16:05:17 Producer ID iqeaoeywrw pushed twitter id:1579198649024000000
    2022/10/16 16:05:17 Producer ID iqeaoeywrw pushed twitter id:1579198600562647040
    2022/10/16 16:05:22 in consumer cgriwwietx redis length:81
    2022/10/16 16:05:22 cgriwwietx consumed  id:1579198600562647040
    2022/10/16 16:05:22 cgriwwietx consumed  id:1579198649024000000
    2022/10/16 16:05:22 cgriwwietx consumed  id:1579199420012580864
    ...
    2022/10/16 16:05:22 cgriwwietx consumed  id:1581630095462400000
    2022/10/16 16:05:22 cgriwwietx consumed  id:1581644123396018177
    ^C2022/10/16 16:05:28 Received signal interrupt for Consumer cgriwwietx, finishing...
    2022/10/16 16:05:28 Received signal interrupt for ProducerService iqeaoeywrw, finishing...
    2022/10/16 16:05:28 Main function done
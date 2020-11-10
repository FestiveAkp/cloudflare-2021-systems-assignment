Compile with go build, run with ./systems-assignment


Testing my profile page:

Requests sent: 15
Fastest response time: 44.7612ms
Slowest response time: 113.3475ms
Mean response time: 82.902167ms
Median response time: 80.8818ms
Successful response ratio: 100%
Error codes encountered: []
Smallest response: 3277 bytes
Largest response: 3277 bytes

Testing google.com:
Requests sent: 15
Fastest response time: 152.2805ms
Slowest response time: 195.2791ms
Mean response time: 162.68478ms
Median response time: 158.694ms
Successful response ratio: 100%
Error codes encountered: []
Smallest response: 48493 bytes
Largest response: 49344 bytes

It seems like perhaps some kind of caching occurs towards the end of the profiler, because
the fastest response time (at least for my CloudFlare site) is much faster than the 
slowest response time. Either way, the average response time for my workers site is pretty good 
compared to the results for google.com.

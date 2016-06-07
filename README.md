# Stressenginio is a stress tool based on the Boom to perform load-test for engine-io/websocket server 
    - Sunny Zh
Usage Example:
    $ boom -c 50000 -ei "http://221.228.83.182:18301"

About boom
Boom is a tiny program that sends some load to a web application. It's similar to Apache Bench ([ab](http://httpd.apache.org/docs/2.2/programs/ab.html)), but with better availability across different platforms and a less troubling installation experience.

Boom is originally written by Tarek Ziade in Python and is available on [tarekziade/boom](https://github.com/tarekziade/boom). But, due to its dependency requirements and my personal annoyance of maintaining concurrent programs in Python, I decided to rewrite it in Go: https://github.com/rakyll/boom




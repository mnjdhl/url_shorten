# URL Shortening Service;written in Golang
#The service receives HTTP Request at port 9988 on end-point "/shorturl".
#The Query Parameter name used is "longURL". 
#e.g. Client can send request as below:
curl http://127.0.0.1:9988/shorturl?longURL=www.efghi.com/longlongurl4nhgjkkkssslkkhhaa

#To build:
go build

#Unit Test:
go test

#More tests
1. Start the url shortening service from one terminal:
./url_shorten

2. Run the curl client from another terminal:
sh url_tests.sh

3. To launch the container, execute script 'launch_container.sh' from the dir where container image 'UrlShorten.tar.gz' is located

CGO_ENABLED=0 go build
mv url_shorten docker
cd docker
rm UrlShorten.tar.gz
sudo docker build -t urlshorten:0.1.0 . -f Dockerfile --rm --no-cache
sudo docker save -o UrlShorten.tar urlshorten:0.1.0
sudo gzip UrlShorten.tar
rm url_shorten
cd -


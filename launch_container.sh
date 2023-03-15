sudo docker stop UrlShorten
sudo docker rm UrlShorten
sudo docker load -i docker/UrlShorten.tar.gz
if [ ! -d  /opt/data/urlshorten ]; then
	sudo mkdir -p /opt/data/urlshorten
fi
sudo docker run -d -it --name UrlShorten --network host -p 9988:9988 -v /opt/data/urlshorten:/data/ urlshorten:0.1.0

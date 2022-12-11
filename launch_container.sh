sudo docker stop UrlShorten
sudo docker rm UrlShorten
sudo docker load -i UrlShorten.tar.gz
sudo docker run -d -it --name UrlShorten --network host -p 9988:9988 -v /opt/data/urlshorten:/data/ urlshorten:0.1.0

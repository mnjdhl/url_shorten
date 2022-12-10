FROM golang:1.16-alpine

WORKDIR /ShortUrl
COPY --chown=0:0 url_shorten ./

RUN chmod 0775 /ShortUrl/url_shorten

#RUN apk add --no-cache libc6-compat

EXPOSE 9988

CMD [ "/ShortUrl/url_shorten" ]


FROM alpine:3.19

RUN apk add --no-cache bash curl jq ca-certificates

COPY chain/build/casedigestord-linux /usr/local/bin/casedigestord
RUN chmod +x /usr/local/bin/casedigestord

EXPOSE 26656 26657

ENTRYPOINT ["casedigestord"]
CMD ["start"]

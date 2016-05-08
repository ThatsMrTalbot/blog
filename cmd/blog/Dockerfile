FROM iron/base

RUN mkdir /blog.git

EXPOSE 80 80

ADD blog /
VOLUME ["/blog.git"]

ENTRYPOINT ["/blog", "-path=/blog.git", "-http=:80"]

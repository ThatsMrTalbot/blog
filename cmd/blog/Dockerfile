FROM iron/base

RUN mkdir /blog.git

EXPOSE 80 80

ADD blog /
RUN chmod +x /blog

VOLUME ["/blog.git"]

CMD /blog -path=/blog.git -http=:80

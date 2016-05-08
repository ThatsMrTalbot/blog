FROM iron/base

RUN mkdir /app
RUN mkdir /blog.git

EXPOSE 80 80

WORKDIR /app
ADD blog /app/
VOLUME ["/blog.git"]

ENTRYPOINT ["./blog", "-path=/blog.git", "-http=:80"]

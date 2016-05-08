FROM iron/base

RUN mkdir /app
RUN mkdir /blog.git


WORKDIR /app
COPY blog /app/
VOLUME ["/blog.git"]
ENTRYPOINT ["./blog -path=/blog.git"]

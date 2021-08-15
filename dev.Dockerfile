FROM golang
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - && apt install -y nodejs libglib2.0
WORKDIR /build

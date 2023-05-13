ARG VERSION
FROM lvsoso/base:${VERSION}

RUN git clone https://github.com/pjd/pjdfstest.git
RUN cd pjdfstest/ && autoreconf -ifs && ./configure && make pjdfstest


RUN wget https://mirrors.edge.kernel.org/pub/software/scm/git/git-2.40.1.tar.gz
RUN tar -zxf  git-2.40.1.tar.gz   \
            && cd git-2.40.1   \
            && make configure \
            && ./configure --prefix=/usr NO_TCLTK=true \
            && make all


RUN apt-get update && apt-get install -y \
    curl \
    && rm -rf /var/lib/apt/lists/*

ADD goup-install.sh .

RUN  chmod +x  ./goup-install.sh && sh ./goup-install.sh  --skip-prompt

RUN .  $HOME/.go/env && goup install 1.19.2 && goup set 1.19.2

RUN git clone https://github.com/go-gitea/gitea && cd gitea && git checkout v1.19.0

RUN .  $HOME/.go/env &&  cd /code/gitea &&  go test  -race  -tags='sqlite sqlite_unlock_notify' code.gitea.io/gitea/modules/git

CMD ["/bin/bash"]
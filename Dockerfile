FROM ruby:2.3.3-slim
MAINTAINER Daisuke Fujita <dtanshi45@gmail.com> (@dtan4)

RUN bundle config --global frozen 1
RUn apt-get update && \
    apt-get install -y cmake && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY Gemfile /app/
COPY Gemfile.lock /app/
RUN bundle install --without test development --system -j4

COPY . /app

EXPOSE 9292

CMD ["bundle", "exec", "rackup", "-E", "production", "-p", "9292"]

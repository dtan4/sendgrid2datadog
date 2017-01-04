require "json"
require "logger"
require "sinatra/base"
require "datadog/statsd"

class App < Sinatra::Base
  configure do
    set :statsd, Datadog::Statsd.new(ENV["DOGSTATSD_HOST"] || "localhost", 8125)
    use Rack::CommonLogger, Logger.new(STDOUT)
  end

  def statsd
    settings.statsd
  end

  get "/" do
    "https://github.com/dtan4/sendgrid2datadog"
  end

  post "/webhook" do
    events = JSON.parse(request.body.read)

    statsd.batch do |s|
      events.each do |event|
        s.increment("sendgrid.event.#{event['event']}")
      end
    end

    "Events was sent to Datadog"
  end
end

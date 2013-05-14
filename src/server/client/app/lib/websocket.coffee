ReliableSocket = require 'lib/ws'
module.exports =
  create: (url) ->
    ws = new ReliableSocket("ws://#{ location.host }#{ url }")
    ws.onerror = (e) ->
      console.log e
    ws
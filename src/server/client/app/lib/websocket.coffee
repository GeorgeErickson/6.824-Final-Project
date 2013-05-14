module.exports =
  create: (url) ->
    ws = new WebSocket("ws://#{ location.host }#{ url }")
    ws.onerror = (e) ->
      console.log e
    ws
websocket = require 'lib/websocket'

module.exports = class DocumentModel extends Backbone.Model
  idAttribute: 'Name'

  send: (data) =>
    ws = @getSocket()
    ws.send_json data
  
  getSocket: =>
    # Caches access to websocket
    unless @_ws
      @_ws = websocket.create "/documents/#{ @get 'Name' }"
      @_ws.onmessage = (e) =>
        @trigger 'message', e.data
      @_ws.onclose = =>
        @_ws = null
    @_ws

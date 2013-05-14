websocket = require 'lib/websocket'

module.exports = class DocumentModel extends Backbone.Model
  idAttribute: 'Name'
  initialize: ->
    ws = @getSocket()
    count = 0
    ws.onmessage = (e) =>
      data = JSON.parse e.data
      data.count = count
      if data.Snapshot
        @set 'Snapshot', data.Snapshot
      @trigger 'message', data
      count += 1

  send: (data) =>
    ws = @getSocket()
    ws.send_json data
  
  getSocket: =>
    # Caches access to websocket
    unless @_ws
      @_ws = websocket.create "/documents/#{ @get 'Name' }"
    @_ws

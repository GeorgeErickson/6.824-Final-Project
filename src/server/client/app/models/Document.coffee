websocket = require 'lib/websocket'

module.exports = class DocumentModel extends Backbone.Model
  idAttribute: 'Name'
  initialize: ->
    ws = @getSocket()
    ws.onmessage = (e) =>
      @trigger 'message', JSON.parse e.data

  send: (data) =>
    ws = @getSocket()
    ws.send_json data
  
  getSocket: =>
    # Caches access to websocket
    unless @_ws
      @_ws = websocket.create "/documents/#{ @get 'Name' }"
    @_ws

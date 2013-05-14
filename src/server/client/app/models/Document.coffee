websocket = require 'lib/websocket'
text = require 'lib/text'

module.exports = class DocumentModel extends Backbone.Model
  idAttribute: 'Name'
  version: null
  ops:
    inflight: null
    pending: null



  queue: (op) =>
    unless _.isArray op
      op = [op]
    
    if @ops.pending == null
      @ops.pending = op
    else
      @ops.pending = text.compose @ops.pending, op
    setTimeout @flush, 0

  flush: =>
    if @ops.inflight or @ops.pending == null
      return

    @ops.inflight = @ops.pending
    @ops.pending = null

    msg = _.clone @attributes
    msg.Version = @version
    msg.OpData = [@ops.inflight]
    @send msg


  send: (data) =>
    ws = @getSocket()
    ws.send_json data
  
  getSocket: =>
    # Caches access to websocket
    unless @_ws
      @_ws = websocket.create "/documents/#{ @get 'Name' }"
      @_ws.onconfirm = (e) =>
        console.log 'Confirm Message sent'
        @version += 1
        @ops.inflight = null
      @_ws.onmessage = (e) =>
        #Ignore Messages We have already seen
        if @version and @version >= e.data.version
          console.log 'This message has already been seen'
          return
        @version = e.data.version

        @trigger 'message', e.data
      @_ws.onclose = =>
        @_ws = null
    @_ws

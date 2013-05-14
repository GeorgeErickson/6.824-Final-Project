websocket = require 'lib/websocket'
text = require 'lib/text'

module.exports = class DocumentModel extends Backbone.Model
  idAttribute: 'Name'
  version: 0
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
      console.log @ops
      return

    @ops.inflight = @ops.pending
    @ops.pending = null

    msg = _.clone @attributes
    msg.Version = @version
    msg.OpData = [@ops.inflight]
    console.log msg
    @send msg


  send: (data) =>
    ws = @getSocket()
    ws.send_json data
  
  getSocket: =>
    # Caches access to websocket
    unless @_ws
      @_ws = websocket.create "/documents/#{ @get 'Name' }"
      @_ws.onconfirm = (e) =>
        console.log 'confirmed'
        @version += 1
        @ops.inflight = null
        @flush()
      @_ws.onmessage = (e) =>
        #Ignore Messages We have already seen
        if @version and @version >= e.data.version
          console.log 'This message has already been seen'
          return
        @version = e.data.Version

        @trigger 'message', e.data
      @_ws.onclose = =>
        @ops.inflight = null
        @ops.pending = null
        @_ws = null
    @_ws

module.exports = class ChatWindowView extends Backbone.Marionette.ItemView
  className: 'ChatWindow'
  template: require 'views/templates/ChatWindow'

  intitialize: (options) ->
    @ws = options.ws
    @ws.onmessage = @onMessage

  events:
    'keypress #chatinput': 'sendKey'
    'click .chatinputbtn': 'sendBtn'

  sendKey: (e) ->
    return  if e.which isnt 13 or not @$('#chatinput').val().trim()
    @send @$('#chatinput').val()
    @$('#chatinput').val ""

  sendBtn: ->
    return if not @$('#chatinput').val().trim()
    @send @$('#chatinput').val()
    @$('#chatinput').val ""

  send: (msg) ->
    @ws.send(msg)

  onMessage: (msg) ->
    @$('.chatFrame').append('<p>'+msg.data+'</p>')
    @$('.chatFrame').scrollTop @$('.chatFrame').scrollHeight

  close: ->
    @ws.close()

  open: (name) ->
    @ws = websocket.create "/chat/#{ name }"
    @ws.onmessage = @onMessage
websocket = require 'lib/websocket'

class Chat extends Backbone.View
  className: 'hide'

  events:
    'keypress #chatinput': 'sendKey'
    'click .chatinputbtn': 'sendBtn'

  show: (model) ->
    @setElement $("#chat")
    @$el.removeClass 'hide'
    @open(model.get('Name'))

  hide: ->
    @setElement $("#chat")
    @$el.addClass 'hide'
    @close()

  sendKey: (e) ->
    return  if e.which isnt 13 or not @$('#chatinput').val().trim()
    @send @$('#chatinput').val()
    @$('#chatinput').val ""

  sendBtn: ->
    return if not @$('#chatinput').val().trim()
    @send @$('#chatinput').val()
    @$('#chatinput').val ""

  send: (msg) ->
    @last_message = msg
    @ws.send_json msg

  onMessage: (msg) =>
    if msg.data == @last_message
      $('.chatFrame').append('<div class="bubble bubble--alt">'+msg.data+'</div>')
      $('.chatFrame').scrollTop $('.chatFrame')[0].scrollHeight
    else
      $('.chatFrame').append('<div class="bubble">'+msg.data+'</div>')
      $('.chatFrame').scrollTop $('.chatFrame')[0].scrollHeight

  close: ->
    @$('.chatFrame').html("")
    @$('#chatinput').val ""
    if @ws?
      @ws.close()
      @last_message = ""

  open: (name) ->
    @ws = websocket.create "/chat/#{ name }"
    @ws.onmessage = @onMessage

chat = new Chat()
module.exports = chat

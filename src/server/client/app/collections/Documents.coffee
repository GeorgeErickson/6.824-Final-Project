websocket = require 'lib/websocket'

class Documents extends Backbone.Collection
  sockets: {}
  chats: {}
  model: require 'models/Document'
  url: '/rest/documents'
  initialize: ->
    @ws = websocket.create '/ws'
    @ws.onmessage = (e) =>
      console.log e

  parse: (response, options) ->
    _.values response

  create: (model, options) =>
    ws = websocket.create "/documents/#{ model.id }"
    cws = websocket.create "/chat/#{ model.id }"
    @sockets[model.id] = ws
    @chats[model.id] = cws
    ws.onmessage = (e) =>
      @add JSON.parse e.data
    cws.onmessage = (e) =>
      


documents = new Documents()
documents.fetch()
module.exports = documents

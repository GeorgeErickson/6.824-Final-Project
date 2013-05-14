websocket = require 'lib/websocket'

class Documents extends Backbone.Collection
  sockets: {}
  chats: {}
  model: require 'models/Document'
  url: '/rest/documents'
  initialize: ->
    @ws = websocket.create '/ws'
    @ws.onmessage = (e) =>
      @add e.data

  parse: (response, options) ->
    _.values response

  create: (model, options) =>
    model = @_prepareModel model, options
    cws = websocket.create "/chat/#{ model.id }"
    @chats[model.id] = cws
    model.once 'message', (data) =>
      console.log data.Name
      @add data
    cws.onmessage = (e) =>
      


documents = new Documents()
documents.fetch()
module.exports = documents

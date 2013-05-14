websocket = require 'lib/websocket'

class Documents extends Backbone.Collection
  sockets: {}
  model: require 'models/Document'
  url: '/rest/documents'
  initialize: ->
    console.log 'dddd'
    @ws = websocket.create '/ws'
    @ws.onmessage = (e) =>
      console.log e

  parse: (response, options) ->
    _.values response

  create: (model, options) =>
    ws = websocket.create "/documents/#{ model.id }"
    @sockets[model.id] = ws

    ws.onmessage = (e) =>
      @add JSON.parse e.data


documents = new Documents()
documents.fetch()
module.exports = documents

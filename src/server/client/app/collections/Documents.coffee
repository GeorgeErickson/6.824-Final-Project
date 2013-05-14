class Documents extends Backbone.Collection
  sockets: {}
  model: require 'models/Document'
  url: '/rest/documents'
  
  parse: (response, options) ->
    _.values response

  create: (model, options) =>
    ws = new WebSocket("ws://#{ location.host }/documents/#{ model.id }")
    @sockets[model.id] = ws

    ws.onmessage = (e) ->
      console.log e


documents = new Documents()
documents.fetch()
module.exports = documents

websocket = require 'lib/websocket'

class Documents extends Backbone.Collection
  sockets: {}
  chats: {}
  model: require 'models/Document'
  url: '/rest/documents'
  initialize: ->
    @ws = websocket.create '/ws'
    @ws.onmessage = (e) =>
      unless @get e.data.Name
        @add e.data


  parse: (response, options) ->
    _.values response

  create: (model, options) =>
    model = new @model model
    
      


documents = new Documents()
documents.fetch()
module.exports = documents

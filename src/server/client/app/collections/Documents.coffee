websocket = require 'lib/websocket'

class Documents extends Backbone.Collection
  sockets: {}
  chats: {}
  model: require 'models/Document'
  url: '/rest/documents'
  initialize: ->
    @ws = websocket.create '/ws'
    @ws.onmessage = (e) =>
      doc = @get e.data.Name
      unless doc
        @add e.data
      

      if e.data.Title != doc.get 'Title'
        doc.set 'Title', e.data.Title

  parse: (response, options) ->
    _.values response

  create: (model, options) =>
    model = new @model model
    
      


documents = new Documents()
documents.fetch()
module.exports = documents

define ['views/home', 'views/document', 'collections/documents'], (HomeView, DocumentView, DocumentsCollection)->
  Router = Backbone.Router.extend
    routes:
      "document/:doc_id/": 'document'
      "": "home"

    home: ->
      HomeView.render()
      app.setView HomeView

    document: (doc_id) ->
      DocumentView.setModel DocumentsCollection.get doc_id
      app.setView DocumentView
      DocumentView.set_editor ace.edit "edit"
  

  new Router()
  
   


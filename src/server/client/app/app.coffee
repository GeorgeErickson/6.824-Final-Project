documents = require 'collections/Documents'
app = new Backbone.Marionette.Application()
loading = require 'views/Loading'

class MainRouter extends Backbone.Router
  routes:
    'edit/:document_id/': 'edit'
    'about': 'about'
    '': 'home'

  _set_nav: (route) ->
    $el = $(app.nav.el)
    $el.find("[data-route]").removeClass 'active'
    $el.find("[data-route=#{ route }]").addClass 'active'
  
  _show_edit: (model) ->
    EditNav = require 'views/EditNav'
    DocumentEdit = require 'views/DocumentEdit'
    app.content.show new DocumentEdit
      model: model
    app.nav.show new EditNav
      model: model
    
  home: ->
    app.nav.show require 'views/DefaultNav'
    app.content.show require 'views/DocumentList'

  edit: (doc_id) ->
    
    model = documents.get(doc_id)
    if model
      @_show_edit model
    else
      app.content.show loading
      documents.once 'sync', =>
        model = documents.get(doc_id)
        @_show_edit model




    

  about: ->
    app.nav.show require 'views/DefaultNav'
    app.content.show require 'views/About'

app.router = new MainRouter()

app.addRegions
  nav: '#nav'
  content: '#content'

app.addInitializer (options) ->
  Backbone.history.start()

module.exports = app
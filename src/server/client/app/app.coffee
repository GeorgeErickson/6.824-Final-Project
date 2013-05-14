documents = require 'collections/Documents'
app = new Backbone.Marionette.Application()
loading = require 'views/Loading'
editor = require 'lib/editor'
chat = require 'lib/chat'

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
    model.getSocket()
    app.content.reset()
    EditNav = require 'views/EditNav'
    editor.show model
    chat.show model
    app.nav.show new EditNav
      model: model
    
  home: ->
    editor.hide()
    chat.hide()
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
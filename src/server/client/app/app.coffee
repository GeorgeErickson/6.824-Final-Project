documents = require 'collections/Documents'
app = new Backbone.Marionette.Application()

class MainRouter extends Backbone.Router
  routes:
    'edit/:document_id/': 'edit'
    'about': 'about'
    '': 'home'

  _set_nav: (route) ->
    $el = $(app.nav.el)
    $el.find("[data-route]").removeClass 'active'
    $el.find("[data-route=#{ route }]").addClass 'active'
    
  home: ->
    app.nav.show require 'views/DefaultNav'
    app.content.show require 'views/DocumentList'

  edit: (doc_id) ->
    EditNav = require 'views/EditNav'
    DocumentEdit = require 'views/DocumentEdit'
    model = documents.get(doc_id)
    app.content.show new DocumentEdit
      model: model
    app.nav.show new EditNav
      model: model

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
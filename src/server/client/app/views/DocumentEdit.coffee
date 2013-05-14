module.exports = class DocumentEditView extends Backbone.Marionette.ItemView
  className: 'DocumentEdit'
  template: require 'views/templates/DocumentEdit'

  events:
    'keyup #editor': 'content_change'

  content_change: (e) ->
    key = e.which
    console.log key
app = require 'app'

module.exports = class DocumentListItemView extends Backbone.Marionette.ItemView
  className: 'DocumentListItem'
  tagName: 'tr'
  template: require 'views/templates/DocumentListItem'

  events:
    'click': 'open_edit'

  open_edit: ->
    app.router.navigate "/edit/#{ @model.id }/", true


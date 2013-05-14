class DocumentList extends Backbone.Marionette.CompositeView
  className: 'DocumentList'
  collection: require 'collections/Documents'
  template: require 'views/templates/DocumentList'
  itemView: require 'views/DocumentListItem'
  itemViewContainer: 'tbody'

  events:
    "click [data-action=create-new]": 'create_new_document'

  create_new_document: ->
    model = @collection.create
      id: uuid.v4()

  initialize: ->
    @render()

module.exports = new DocumentList()

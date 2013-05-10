#Home
define ['collections/documents', 'views/home-item'], (DocumentsCollection, HomeItemView) ->
  Thorax.View.extend 
    name: 'home'
    collection: DocumentsCollection

  new Thorax.Views.home()
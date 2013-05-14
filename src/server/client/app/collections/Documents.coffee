class Documents extends Backbone.Collection
  model: require 'models/Document'
  localStorage: new Backbone.LocalStorage("DocumentsCollection")

documents = new Documents()
documents.fetch()
module.exports = documents

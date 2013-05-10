define [], ->
  class Documents extends Backbone.Collection



  unless Thorax.Collections.Documents
     Thorax.Collections.Documents = new Documents([{id:1, name:'Room 1'}, {id:2, name:'Room 2'}])

  Thorax.Collections.Documents
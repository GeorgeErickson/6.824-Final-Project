#Home
define [], ->
  
  Thorax.View.extend 
    name: 'homeitem'
    template: 'home-item'
    tagName: 'tr'
    events:
      'click': 'detail'

    detail: ->
      app.router.navigate "//document/#{ @model.id }/", 
        trigger: true
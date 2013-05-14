DefaultNavView = require 'views/DefaultNav'

describe 'DefaultNavView', ->
    beforeEach ->
        @view = new DefaultNavView()

    it 'should exist', ->
        expect(@view).to.be.ok

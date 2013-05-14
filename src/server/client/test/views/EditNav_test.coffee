EditNavView = require 'views/EditNav'

describe 'EditNavView', ->
    beforeEach ->
        @view = new EditNavView()

    it 'should exist', ->
        expect(@view).to.be.ok

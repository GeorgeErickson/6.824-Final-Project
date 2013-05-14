AboutView = require 'views/About'

describe 'AboutView', ->
    beforeEach ->
        @view = new AboutView()

    it 'should exist', ->
        expect(@view).to.be.ok

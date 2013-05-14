exports.config =
  files:
    javascripts:
      joinTo:
        'scripts/app.js': /^app/
        'scripts/vendor.js': /^vendor/
      order:
        before: [
          'vendor/scripts/console-polyfill.js'
          'vendor/scripts/jquery.js'
          'vendor/scripts/underscore.js'
          'vendor/scripts/backbone.js'
          'vendor/scripts/bootstrap.min.js'
          'vendor/scripts/bootstrap-editable.js'
        ]

    stylesheets:
      joinTo:
        'stylesheets/app.css'
      order:
        before: [
          'vendor/styles/normalize.css'
          'vendor/styles/typeplate-unminified.css'
        ]

    templates:
      joinTo: 'scripts/app.js'

  plugins:
    coffeelint:
      pattern: /^app\/.*\.coffee$/

      options:
        indentation:
          value: 2
          level: "error"

        max_line_length:
          value: 180
          level: "error"

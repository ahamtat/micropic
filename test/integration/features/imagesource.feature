
Feature: Image source
  As HTTP proxy for multiple microservices
  In order to understand that previewer receives and processes images correctly
  I want to receive response and check it

  Scenario: Image source is not available
    When I make a "GET" request to "http://proxy:8080/fill/50/50/imagesource.not/_gopher_original_1024x504.jpg"
    Then I get response status code 503
    And the response body is "Service Unavailable"

  Scenario: Image not found
    When I make a "GET" request to "http://proxy:8080/fill/50/50/imagesource/_gopher_not_found.jpg"
    Then I get response status code 404

  Scenario: Wrong image
    When I make a "GET" request to "http://proxy:8080/fill/50/50/imagesource/_gopher_not_image.jpg"
    Then I get response status code 500
    And the response body is "image: unknown format"

  Scenario: Image source returns error
    When I make a "GET" request to "http://proxy:8080/fill/50/50/imagesource/_gopher_not_found.jpg"
    Then I get response status code 404

  Scenario: Image source returns image
    When I make a "GET" request to "http://proxy:8080/fill/500/500/imagesource/_gopher_original_1024x504.jpg"
    Then I get response status code 200
#    And the response headers has:
#      | Test-Previewlocation | previewer |

  Scenario: Preview size is correct
    When I make a "GET" request to "http://proxy:8080/fill/50/50/imagesource/_gopher_original_1024x504.jpg"
    Then I get response status code 200
    And preview size is:
      | Width  | 50 |
      | Height | 50 |

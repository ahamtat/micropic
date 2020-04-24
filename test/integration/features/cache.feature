
Feature: Cache
  As HTTP proxy for multiple microservices
  In order to understand that cache saves and returns previews
  I want to receive response and check it

  Scenario: Cache returns preview
    Given I make a "GET" request to "http://localhost:8080/fill/200/300/imagesource/_gopher_original_1024x504.jpg"
    When I make a "GET" request to "http://localhost:8080/fill/200/300/imagesource/_gopher_original_1024x504.jpg"
    Then I get response status code 200
    And the response headers has:
      | Test-Previewlocation | cache |
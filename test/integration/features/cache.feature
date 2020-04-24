
Feature: Cache
  As HTTP proxy for multiple microservices
  In order to understand that cache saves and returns previews
  I want to receive response and check it

  Scenario: Cache returns preview
    Given I make a "GET" request to "http://localhost:8080/fill/50/50/imagesource/_gopher_original_1024x504.jpg"
    When I make a "GET" request to "http://localhost:8080/fill/50/50/imagesource/_gopher_original_1024x504.jpg"
    Then I get a "200" response
    And the proxy response has details:
      | Header | Test-Previewlocation | cache |
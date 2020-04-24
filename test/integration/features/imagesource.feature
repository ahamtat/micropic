
Feature: Image source
  As HTTP proxy for multiple microservices
  In order to understand that previewer receives and processes images correctly
  I want to receive response and check it

  Scenario: Image source is not available
    When I make a "GET" request to "http://localhost:8080/fill/50/50/imagesource.not/_gopher_original_1024x504.jpg"
    Then I get a "503" response
    And the proxy response has details:
      | Body | Service Unavailable |

  Scenario: Image not found
    When I make a "GET" request to "http://localhost:8080/fill/50/50/imagesource/_gopher_not_found.jpg"
    Then I get a "404" response

  Scenario: Wrong image
    When I make a "GET" request to "http://localhost:8080/fill/50/50/imagesource/_gopher_not_image.jpg"
    Then I get a "500" response
    And the proxy response has details:
      | Body | image: unknown format |

  Scenario: Image source returns error
    When I make a "GET" request to "http://localhost:8080/fill/50/50/imagesource/_gopher_not_found.jpg"
    Then I get a "404" response

  Scenario: Image source returns image
    When I make a "GET" request to "http://localhost:8080/fill/50/50/imagesource/_gopher_original_1024x504.jpg"
    Then I get a "200" response
    And the proxy response has details:
      | Header | Test-Previewlocation | previewer |

  Scenario: Preview size is correct
    When I make a "GET" request to "http://localhost:8080/fill/50/50/imagesource/_gopher_original_1024x504.jpg"
    Then I get a "200" response
    And preview size is:
      | Width  | 50 |
      | Height | 50 |

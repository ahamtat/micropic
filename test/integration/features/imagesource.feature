
Feature: Image source
  As HTTP proxy for multiple microservices
  In order to understand that previewer receives and processes images correctly
  I want to receive response and check its status

  Scenario: Image source is not available
    Given Connection to Calendar API on "api:8888"
    And There is the event:
    """
		{
			"title": "Event 1",
			"description": "Data for testing microservices",
			"owner": "Artem",
			"startTime": "2020-04-02T12:03:00+03:00"
		}
	"""
    When I send AddEvent request
    Then response should have event
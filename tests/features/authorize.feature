Feature: authorize
  In order to authorize a request
  As a user trying to access the resources
  I need to be able to authorize the request with a token

  Scenario: authorize a request
    Given I have a valid token
    When I authorize the request
    Then the request should be authorized
# Summary

Language: GO

* Features:
  * Player moves
    * Deploy troops
    * Move troops

* Core functionality
  * Start new game based on configured settings
    * Default - load initial game state from file
  * Validate moves sent by players
  * Execute player moves after all players submitted their moves
  * Give back game state to Backend after all moves have been executed

* Links
  1. Engine <=> Backend : HTTP (2 way communication)

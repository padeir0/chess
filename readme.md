# Simplified Chess

Implements a command line UI for simplified chess and a chess engine.

## Rules

Rules differing from standard chess:

 - No check or checkmate, the game is won by capturing the king
 - King can put himself into check (possibly losing the game)
 - Passing your turn is allowed
 - If both players pass turns successively, the game ends in a draw
 - No castling
 - No en passant
 - Pawns move only a single square
 - No repetition rules
 - Pawns promote automatically to queens
 - 50 moves without a capture ends in a draw

The rest is the same as classical chess (i think).

## Commands

```
move a2 a4   // moves piece at a2 to a4
move a7 a8 q // moves piece at a7 to a8 and specifies promotion

save mypoint    // saves this current position as "mypoint"
restore mypoint // restores board position to "mypoint"
restore         // restores board to initial position

show            // shows board
show moves      // shows valid moves

profile <label>
stopprofile

quit         // quits
exit         // quits
clear        // clears screen
```

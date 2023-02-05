# Chess

Implements a command line UI for chess and
a simple mailbox chess engine.

The goal is to have a 1500+ elo engine, eventually implementing UCI,
both for the engine and the UI.

Currently it does only marginally better than random moves.

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

NO // makes the engine cry
```

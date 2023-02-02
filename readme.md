# Chess

Implements a command line UI for chess and
a simple mailbox chess engine.
The goal is to have a 1500+ elo engine, eventually implementing UCI,
both for the engine and the UI.

## Commands

`next`, `undo`, `show attacked/defended/vulnerable` are not implemented

```
next 3       // returns top 3 best next moves (for the current player)
next 3 1     // returns top 3 best next moves (for the current player)
             // and the top 1 best response for each of them

move a2 a4   // moves piece at a2 to a4
move a7 a8 q // moves piece at a7 to a8 and specifies promotion

undo         // undos the previous move
undo 3       // undos the 3 previous moves

save mypoint    // saves this current position as "mypoint"
restore mypoint // restores board position to "mypoint"
restore         // restores board to initial position

show            // shows board
show attacked   // shows board and colors pieces being attacked
show defended   // shows board and colors pieces being defended
show vulnerable // shows board and colors vulnerable pieces

quit         // quits
exit         // quits
clear        // clears screen
```

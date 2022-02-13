# Gambit

A client to play chess in your terminal. All this work are inspired and based in the original [gambit](https://github.com/maaslalani/gambit) from [Maas Lalani](https://github.com/maaslalani). I get it and build [gambitsrv](https://github.com/sdemingo/gambitsrv),  a simple server in Go to manage chess games between human players. Gambit is specially designed to work on a unix shared server (*aka tilde or similar*)

<br/>
<p align="center">
  <img width="90%" src="./chess.gif?raw=true" alt="Terminal chess" />
</p>
<br/>


### Move

There are two ways to move in `gambit`:

* Type out the square the piece you want to move is on, then type out the square to which you want to move the piece.
* With the mouse, click on the target piece and target square.

### Play

You can create a new a game by running:

```
gambit 
```

Then, you get the game id from the server. You should comunicate it to the other player. If the game has been created by other player, you can join with:
```
gambit -g <id-of-game>
```

You can press <kbd>ctrl+f</kbd> to flip the board to give a better perspective
for the second player.


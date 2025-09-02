# Punto Banco (Baccarat) on Go

A Go implementation of the _punto banco_ card game (a version of [baccarat](https://en.wikipedia.org/wiki/Baccarat)), featuring a terminal-based UI built with the [Bubble Tea framework](https://github.com/charmbracelet/bubbletea/).

> This implementation does not include a gambling component — no real bets or payouts are involved in the game.

## Running the Application

Execute the binary or run the application directly with:

```bash
go run cmd/main.go
```

<img width="530" src="./screenshot.png" />

### Features

- 6 decks in the shoe
- Casino-style shuffling with shoe cutting and card burning
- Infinity game (shoe updates automatically when it ends)
- Game session statistics
- Terminal-based UI

## Game Rules

_Punto banco_ is a simplified version of baccarat where each move is determined by the drawn cards. The game proceeds according to fixed rules, and players make no decisions during the coup (round).

The game uses standard baccarat card values:

- Ace = 1
- 2–9 = pip value
- 10 and face cards (Jack, Queen, and King) = 0

The objective is to predict which of the two hands, the Punto (player) or the Banco (banker), will have a total closest to nine.

Player bets on either:

- Punto (player)
- Banco (banker)
- Egalité (tie)

The dealer deals two cards to both the player and the banker.

**If the total of a hand is 10 or more, only the last digit is counted** (modulo 10).

If either the player or the banker (or both) has a total of **8 or 9**, the round ends immediately; it is referred to as a «natural».

Otherwise:

- If the player's total is **0–5**, they draw a third card.
- If the player's total is **6 or 7**, they stand.

The banker's decision to draw a third card depends on both their total and the player's third card, following a predefined «tableau»:

**Banker's total / Player's third card value**

|         | 0   | 1   | 2   | 3   | 4   | 5   | 6   | 7   | 8   | 9   |
| ------- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| **0-2** | H   | H   | H   | H   | H   | H   | H   | H   | H   | H   |
| **3**   | H   | H   | H   | H   | H   | H   | H   | H   | S   | H   |
| **4**   | S   | S   | H   | H   | H   | H   | H   | H   | S   | S   |
| **5**   | S   | S   | S   | S   | H   | H   | H   | H   | S   | S   |
| **6**   | S   | S   | S   | S   | S   | S   | H   | H   | S   | S   |
| **7+**  | S   | S   | S   | S   | S   | S   | S   | S   | S   | S   |

Legend: H (Hit, draw another card), S (Stand, no more cards).

After all cards are drawn, the **hand with the total closest to nine wins**.

All possible combinations are covered with unit tests.

---

## Simulator

Simulator exists as a separate program:

```bash
go run cmd/simulator/main.go
```

This simulator runs the _punto banco_ game, and during each round, it bets on Punto (player), Banco (banker), or Égalité (tie) depending on the chosen strategy.

«The game» is a game session, in which the **simulation starts with the bankroll of $1000** and ends when it cannot afford to bet the next bet.

**The default bet in the simulator is $10**, because it is the minimum bet for baccarat in Las Vegas. Therefore, for flat betting strategies, the simulator bets 1% of its initial bankroll in each round.

The simulator features a logic-based betting approach: the game ends if the current simulator's bankroll falls below the amount required for the next round.

- In flat betting strategies, the game ends when the bankroll becomes 0.
- In progression strategies, the game ends when the number of consecutive wins becomes too favorable (or too negative) that the simulator needs to bet more than the bankroll allows. In this case, the bankroll can be higher than 0 (even too big for edge cases).

Payouts (or pop-up of the bankroll in the context of the simulator) in simulation are made according to standard baccarat rules:

- **1-to-1** on Punto bets.
- **19-to-20** on Banco bets (5% commission is designed to balance Banco's statistical advantage of a slightly higher probability of winning).
- **8-to-1** on Égalité bet.

Simulator implements the following strategies:

- Bet always on Punto (player)
- Bet always on Banco (banker)
- Bet always on Égalité (tie)
- Bet on last hand
- Bet on random
- Martingale
- Paroli
- Fibonacci
- D'Alembert
- 1-3-2-6

The statistics of simulations include the following items:

- Mean rounds per game session until the moment when the gambler can no longer bet.
- Minimum number of played rounds per game session across all simulations.
- Maximum number of played rounds per game session across all simulations.
- Mean wins per game session.
- Minimum wins per session across all simulations. For progression strategies, if a gambler gets into a series of losses, she may run out of bankroll before the first win, and therefore, the minimum number of wins will be 0.
- Maximum wins per game session across all simulations.
- Win rate — the percentage of rounds that a gambler wins over the number of played rounds. It is the way to measure the effectiveness of a strategy.
- The rate of zero-win games indicates the percentage of game sessions that ended without a single win occurring. It may serve as an indicator of the amount of risk associated with a strategy.
- Mean winning streak.
- Maximum winning streak per game session across all simulations.
- Mean losing streak.
- Maximum losing streak per game session across all simulations.
- Mean peak bankroll per game session.
- Maximum recorded bankroll across all simulations. It is the maximum winning amount that occurred in the simulation session of a chosen strategy.
- Profitable games — the percentage of game sessions with a profit opportunity, or the percentage of game sessions in which the bankroll exceeded 101% of the initial value. It shows the percentage of games in which the gambler hit a profit target (in this case $1010 and above) and could have been in profit (won money) if he had stopped betting.
- Profitably ended games — the percentage of game sessions ended with profit, or the percentage of game sessions that end when the current bankroll exceeds 101% of the initial value. This edge case was explained above.

The data of the simulation of each strategy run on 1M game session simulation is the basis of the article «[When You Run Out of Money Playing Baccarat (Punto Banco)](https://adequatica.github.io/2025/09/02/when-you-run-out-of-money-playing-baccarat-punto-banco.html)».

## Race Results Screen Behavior

- **Date**: 2026-01-09
- **Context**: Align race finish UX with expected two-player flow and ELO system.

### Changes

- **Result timing**: The race results screen in `race/[code]/page.js` now appears **only after all players have finished** (i.e. every participant with role `player` has a final verdict), instead of showing immediately when the first player gets `AC`.
- **Race updates**: The race room listens to WebSocket `race_event` messages:
  - `player_finished`: triggers a fresh race fetch so `verdict` and `finish_time` are up-to-date for all players.
  - `elo_updated`: captures the `elo_changes` payload for that race.
- **Results UI**: The results list now shows, per player:
  - Medal / rank
  - Nickname
  - Final verdict and finish time
  - ELO rating delta (e.g. `+28`, `-15`, or `0` if unchanged)

### Notes

- Spectators continue to use the dedicated `SpectatorView` and do not see the race results screen directly.
- ELO changes are driven by the backend `elo_updated` race event described in `docs/elo-system.md`; only players with `AC` are included in the rating calculation, others display a `0` delta on the client.


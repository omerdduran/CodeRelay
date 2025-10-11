# Verdict WebSocket Contract

The API streams head-to-head match verdicts over a dedicated WebSocket channel. Messages are sent as JSON objects with the following shape:

```json
{
  "type": "match.verdict",
  "matchId": "match_123",
  "payload": {
    "player": {
      "id": "player_a",
      "nickname": "CoderA"
    },
    "submission": {
      "id": "sub_456",
      "language": "python",
      "submittedAt": "2025-01-15T18:04:05Z"
    },
    "verdict": {
      "code": "AC",
      "runtimeMs": 742,
      "durationMs": 12500
    },
    "match": {
      "phase": "in_progress",
      "winnerId": null
    }
  }
}
```

- `type` distinguishes verdict pushes from future event families.
- `matchId` allows clients to ignore events for other rooms if they share a socket.
- `payload.player` describes which contender the verdict belongs to; `id` maps to the `match_participants` table.
- `payload.submission` mirrors the `submissions` database row so the UI can display timestamps and languages without an extra fetch.
- `payload.verdict.code` is one of `AC`, `WA`, `TLE`. A future `RE` (runtime error) fits the same contract.
- `payload.verdict.runtimeMs` captures judge runtime; `durationMs` measures wall time from match start to verdict.
- `payload.match.phase` is `pending`, `in_progress`, or `finished`. When the duel ends, `winnerId` is set and final `match.summary` events can be emitted using the same envelope.

Heartbeat messages reuse the envelope with `type: "meta.heartbeat"` and an empty payload to keep connections alive.

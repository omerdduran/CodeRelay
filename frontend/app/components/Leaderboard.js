'use client';

import { useEffect, useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { useWebSocket } from '../hooks/useWebSocket';
import styles from './Leaderboard.module.css';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// If problemId is null => global ELO leaderboard
// If problemId is set   => per-problem ranking (by solve time)
export default function Leaderboard({ problemId = null }) {
  const { t } = useTranslation();
  const [entries, setEntries] = useState([]);
  const [loading, setLoading] = useState(true);

  const isGlobal = problemId == null;

  const fetchLeaderboard = useCallback(async () => {
    try {
      let url = `${API_URL}/api/leaderboard`;

      if (isGlobal) {
        // Global ELO leaderboard
        url += '?type=elo&limit=100';
      } else {
        // Per-problem leaderboard
        url += `?problem_id=${problemId}`;
      }

      const res = await fetch(url);
      if (res.ok) {
        const data = await res.json();
        setEntries(data || []);
      }
    } catch (err) {
      console.error('Failed to fetch leaderboard:', err);
    } finally {
      setLoading(false);
    }
  }, [isGlobal, problemId]);

  // Handle WebSocket messages
  const handleMessage = useCallback(
    (message) => {
      if (message.type === 'submission_update') {
        const payload = JSON.parse(message.payload);
        if (!isGlobal && payload.status === 'AC' && payload.problem_id === problemId) {
          // Refresh problem leaderboard when someone solves this problem
          fetchLeaderboard();
        }
      } else if (message.type === 'race_event') {
        const payload = JSON.parse(message.payload);
        if (isGlobal && payload.event === 'elo_updated') {
          // Refresh global ELO leaderboard when ratings are updated
          fetchLeaderboard();
        }
      }
    },
    [isGlobal, problemId, fetchLeaderboard],
  );

  const { connected } = useWebSocket(handleMessage);

  useEffect(() => {
    fetchLeaderboard();
  }, [fetchLeaderboard]);

  if (loading) {
    return (
      <div className={styles.container}>
        <div className={styles.loading}>{t('leaderboard.loading')}</div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h3 className={styles.title}>
          {isGlobal ? t('leaderboard.globalRanking') : t('leaderboard.problemRanking')}
        </h3>
        <span className={styles.status}>
          {connected ? t('leaderboard.live') : t('leaderboard.offline')}
        </span>
      </div>

      {entries.length === 0 ? (
        <div className={styles.empty}>{t('leaderboard.empty')}</div>
      ) : (
        <div className={styles.list}>
          {entries.map((entry) => (
            <div key={entry.user_id} className={styles.entry}>
              <span className={styles.rank}>#{entry.rank}</span>
              <span className={styles.nickname}>{entry.nickname}</span>
              {isGlobal ? (
                <span className={styles.elo}>{entry.elo_rating ?? 1200}</span>
              ) : (
                <span className={styles.time}>
                  {entry.solve_time_ms ? `${entry.solve_time_ms}ms` : '-'}
                </span>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
